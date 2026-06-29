package nebula

import (
	"context"
	"crypto/rand"
	"fmt"
	"net/netip"
	"strings"
	"time"

	"github.com/openbao/openbao/sdk/v2/framework"
	"github.com/openbao/openbao/sdk/v2/helper/errutil"
	"github.com/openbao/openbao/sdk/v2/logical"
	"github.com/slackhq/nebula/cert"
	"golang.org/x/crypto/ed25519"
)

func buildPathListCerts(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "certs/",

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: "nebula",
			OperationSuffix: "certs",
		},

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ListOperation: &framework.PathOperation{
				Callback: b.pathCertList,
			},
		},

		HelpSynopsis:    "List all Certificates",
		HelpDescription: "List the fingerprints of all certificates",
	}
}

func (b *backend) pathCertList(ctx context.Context, req *logical.Request, _ *framework.FieldData) (*logical.Response, error) {
	entries, err := req.Storage.List(ctx, "certs/")
	if err != nil {
		return nil, err
	}

	for i, str := range entries {
		entries[i] = formatFingerprint(str)
	}

	caStorageEntry, err := req.Storage.Get(ctx, "ca")
	var cse CertStorageEntry
	caStorageEntry.DecodeJSON(&cse)
	nc, _, _ := cert.UnmarshalCertificateFromPEM([]byte(cse.Pem))

	fingerprint, _ := nc.Fingerprint()
	fingerprintStr := fmt.Sprintf("%v", fingerprint)
	entries = append(entries, formatFingerprint(fingerprintStr))

	return logical.ListResponse(entries), nil
}

func buildPathCert(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "cert/" + framework.MatchAllRegex("fingerprint"),
		Fields: map[string]*framework.FieldSchema{
			"fingerprint": {
				Type:        framework.TypeString,
				Description: `Required: fingerprint of the certificate`,
				Required:    true,
			},
		},
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.pathReadCert,
				Summary:  "",
			},
		},
	}
}

func (b *backend) pathReadCert(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	fingerprint := data.Get("fingerprint").(string)
	if fingerprint == "" {
		return nil, fmt.Errorf("Please Specify Certificate Fingerprint")
	}

	if len(fingerprint) != 79 {
		return nil, fmt.Errorf("Invalid Fingerprint")
	}

	cleanFingerprint := strings.ReplaceAll(fingerprint, ":", "")
	storageEntry, err := req.Storage.Get(ctx, "certs/"+cleanFingerprint)
	if err != nil {
		return nil, fmt.Errorf("Invalid Fingerprint")
	}

	if storageEntry == nil {
		return nil, fmt.Errorf("Certificate not found")
	}

	var cse CertStorageEntry
	storageEntry.DecodeJSON(&cse)
	nc, _, _ := cert.UnmarshalCertificateFromPEM([]byte(cse.Pem))

	pemCert, err := nc.MarshalPEM()

	var ipNetStrings []string
	for _, ipNet := range nc.Networks() {
		ipNetStrings = append(ipNetStrings, ipNet.String())
	}

	resp := &logical.Response{
		Data: map[string]interface{}{
			"notAfter":    nc.NotAfter().Format("02.01.2006 15:04:05"),
			"name":        nc.Name(),
			"ip":          strings.Join(ipNetStrings, ", "),
			"cert":        string(pemCert),
			"fingerprint": formatFingerprint(fingerprint),
		},
	}

	revokedStorageEntry, err := req.Storage.Get(ctx, "revoked/"+cleanFingerprint)
	if err == nil && revokedStorageEntry != nil {
		var revocationDetails RevocationDetails
		err = revokedStorageEntry.DecodeJSON(&revocationDetails)
		if err != nil {
			return nil, fmt.Errorf("failed to decode revocation details: %v", err)
		}

		resp.Data["revocation_time"] = revocationDetails.RevokedAt.Unix()
		resp.Data["revocation_time_rfc3339"] = revocationDetails.RevokedAt.Format(time.RFC3339)
	}

	return resp, nil
}

func buildPathSign(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "sign/" + framework.GenericNameRegex("name"),
		Fields: map[string]*framework.FieldSchema{
			"name": {
				Type:        framework.TypeString,
				Description: `Required: name of the certificate authority`,
				Required:    true,
			},
			"duration": {
				Type:        framework.TypeString,
				Description: `Optional: amount of time the certificate should be valid for. Valid time units are seconds: "s", minutes: "m", hours: "h". Without passing a -duration XXhXXmXXs flag, certificates will be valid up until one second before their signing CA expires.`,
			},
			"groups": {
				Type:        framework.TypeString,
				Description: `Optional: list of groups. This will limit which groups subordinate certs can use.`,
				Default:     "",
			},
			"ip": {
				Type:        framework.TypeString,
				Description: `Required: ipv4 address and network in CIDR notation to assign the cert.`,
				Required:    true,
			},
			"subnets": {
				Type:        framework.TypeString,
				Description: `Optional: list of ip and network in CIDR notation. This will limit which subnet addresses and networks subordinate certs can use.`,
				Default:     "",
			},
		},
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathSign,
				Summary:  "",
			},
		},
	}
}

func (b *backend) pathSign(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	caKeyStorageEntry, err := req.Storage.Get(ctx, "ca_key")
	if err != nil {
		return nil, errutil.InternalError{Err: fmt.Sprintf("unable to fetch nebula ca: %v", err)}
	}

	var caPrivateKey ed25519.PrivateKey
	if err := caKeyStorageEntry.DecodeJSON(&caPrivateKey); err != nil {
		return nil, errutil.InternalError{Err: fmt.Sprintf("unable to decode Nebula CA Key: %v", err)}
	}

	nebulaCACertEntry, err := req.Storage.Get(ctx, "ca")
	if err != nil {
		return nil, errutil.InternalError{Err: fmt.Sprintf("unable to fetch nebula ca: %v", err)}
	}

	var cse CertStorageEntry
	if err := nebulaCACertEntry.DecodeJSON(&cse); err != nil {
		return nil, errutil.InternalError{Err: fmt.Sprintf("unable to decode Nebula Certificate: %v", err)}
	}

	caCert, _, err := cert.UnmarshalCertificateFromPEM([]byte(cse.Pem))
	if err != nil {
		return nil, errutil.InternalError{Err: fmt.Sprintf("unable to parse Nebula Certificate PEM: %v", err)}
	}

	name := data.Get("name").(string)
	if name == "" {
		return nil, fmt.Errorf("nebula Certificate Name may not be empty")
	}

	groups := data.Get("groups").(string)
	_groups := parseGroups(groups)

	duration, durationOk := data.GetOk("duration")

	var _duration time.Duration
	if !durationOk {
		_duration = time.Until(caCert.NotAfter()) - time.Second*1
	} else {
		_duration, err = time.ParseDuration(duration.(string))
		if err != nil {
			return nil, fmt.Errorf("Invalid time format: %s", err)
		}
	}

	ip := data.Get("ip").(string)
	var _ip []netip.Prefix
	if ip != "" {
		rs := strings.Trim(ip, " ")
		if rs != "" {
			prefix, err := netip.ParsePrefix(rs)
			if err != nil {
				return nil, fmt.Errorf("invalid ip definition: %s", err)
			}
			_ip = append(_ip, prefix)
		}
	}

	subnets := data.Get("subnets").(string)
	var _subnets []netip.Prefix
	if subnets != "" {
		for _, rs := range strings.Split(subnets, ",") {
			rs := strings.Trim(rs, " ")
			if rs != "" {
				prefix, err := netip.ParsePrefix(rs)
				if err != nil {
					return nil, fmt.Errorf("invalid subnet definition: %s", err)
				}
				_subnets = append(_subnets, prefix)
			}
		}
	}

	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, errutil.InternalError{Err: fmt.Sprintf("Failed to generate keypair: %v", err)}
	}

	tbs := cert.TBSCertificate{
		Version:        cert.Version2,
		Name:           name,
		Groups:         _groups,
		Networks:       _ip,
		UnsafeNetworks: _subnets,
		NotBefore:      time.Now(),
		NotAfter:       time.Now().Add(_duration),
		PublicKey:      publicKey,
		IsCA:           false,
	}

	newCertificate, err := tbs.Sign(caCert, cert.Curve_CURVE25519, caPrivateKey)
	if err != nil {
		return nil, errutil.InternalError{Err: fmt.Sprintf("failed to sign certificate: %v", err)}
	}

	pemCert, _ := newCertificate.MarshalPEM()
	fingerprint, _ := newCertificate.Fingerprint()
	fingerprintStr := fmt.Sprintf("%v", fingerprint)

	entry, err := logical.StorageEntryJSON("certs/"+fingerprintStr, CertStorageEntry{Pem: string(pemCert)})
	if err != nil {
		return nil, err
	}

	err = req.Storage.Put(ctx, entry)
	if err != nil {
		return nil, err
	}

	resp := &logical.Response{
		Data: map[string]interface{}{
			"notAfter":    time.Now().Add(_duration).Format("02.01.2006 15:04:05"),
			"name":        newCertificate.Name(),
			"cert":        string(pemCert),
			"private_key": string(cert.MarshalPrivateKeyToPEM(cert.Curve_CURVE25519, privateKey)),
			"fingerprint": formatFingerprint(fingerprintStr),
		},
	}

	return resp, err
}
