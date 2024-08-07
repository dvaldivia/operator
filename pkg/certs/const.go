// This file is part of MinIO Operator
// Copyright (c) 2021 MinIO, Inc.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package certs

const (

	// CertsCADir Directory contains all CA certificates other than system defaults for HTTPS.
	CertsCADir = "CAs"

	// PublicCertFile Public certificate file for HTTPS.
	PublicCertFile = "public.crt"

	// TLSCertFile Public certificate file for HTTPS.
	TLSCertFile = "tls.crt"

	// PrivateKeyFile Private key file for HTTPS.
	PrivateKeyFile = "private.key"

	// CAPublicCertFile  Public certificate file for Certificate authority.
	CAPublicCertFile = "ca.crt"

	// TLSKeyFile Private key file for HTTPS.
	TLSKeyFile = "tls.key"
)
