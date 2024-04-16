/*
 * Numbers API
 *
 * The Numbers API enables you to manage your existing numbers and buy new virtual numbers for use with Nexmo's APIs. Further information is here: <https://developer.nexmo.com/numbers/overview>
 *
 * API version: 1.0.18
 * Contact: devrel@nexmo.com
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package number
// Availablenumber struct for Availablenumber
type Availablenumber struct {
	// The two character country code in ISO 3166-1 alpha-2 format
	Country string `json:"country,omitempty"`
	// An available inbound virtual number.
	Msisdn string `json:"msisdn,omitempty"`
	// The type of number: `landline`, `landline-toll-free` or `mobile-lvn`
	Type string `json:"type,omitempty"`
	// The monthly rental cost for this number, in Euros
	Cost string `json:"cost,omitempty"`
	// The capabilities of the number: `SMS` or `VOICE` or `SMS,VOICE` or `SMS,MMS` or `VOICE,MMS` or `SMS,MMS,VOICE`
	Features []string `json:"features,omitempty"`
}
