package models

import (
	"github.com/ottemo/foundation/env"
)

// Package global constants
const (
	ConstErrorModule = "models"
	ConstErrorLevel  = env.ConstErrorLevelModel
)

// StructListItem represents type to hold business layer object information within collection
type StructListItem struct {
	ID    string
	Name  string
	Image string
	Desc  string

	Extra map[string]interface{}
}

// StructAttributeInfo represents type to hold business layer object attribute information
type StructAttributeInfo struct {
	Model      string
	Collection string
	Attribute  string
	Type       string
	Label      string
	IsRequired bool
	IsStatic   bool
	Group      string
	Editors    string
	Options    string
	Default    string
	Validators string
	IsLayered  bool
	IsPublic   bool
}

// Package global variables
var (
	declaredModels = map[string]InterfaceModel{}

	ConstStatesList = map[string]string{
		"":   "No state",
		"AL": "Alabama",
		"AK": "Alaska",
		"AZ": "Arizona",
		"AR": "Arkansas",
		"CA": "California",
		"CO": "Colorado",
		"CT": "Connecticut",
		"DE": "Delaware",
		"FL": "Florida",
		"GA": "Georgia",
		"HI": "Hawaii",
		"ID": "Idaho",
		"IL": "Illinois",
		"IN": "Indiana",
		"IA": "Iowa",
		"KS": "Kansas",
		"KY": "Kentucky",
		"LA": "Louisiana",
		"ME": "Maine",
		"MD": "Maryland",
		"MA": "Massachusetts",
		"MI": "Michigan",
		"MN": "Minnesota",
		"MS": "Mississippi",
		"MO": "Missouri",
		"MT": "Montana",
		"NE": "Nebraska",
		"NV": "Nevada",
		"NH": "New Hampshire",
		"NJ": "New Jersey",
		"NM": "New Mexico",
		"NY": "New York",
		"NC": "North Carolina",
		"ND": "North Dakota",
		"OH": "Ohio",
		"OK": "Oklahoma",
		"OR": "Oregon",
		"PA": "Pennsylvania",
		"RI": "Rhode Island",
		"SC": "South Carolina",
		"SD": "South Dakota",
		"TN": "Tennessee",
		"TX": "Texas",
		"UT": "Utah",
		"VT": "Vermont",
		"VA": "Virginia",
		"WA": "Washington",
		"WV": "West Virginia",
		"WI": "Wisconsin",
		"WY": "Wyoming",
	}

	ConstCountriesList = map[string]string{
		"AF": "Afghanistan",
		"AX": "Åland Islands",
		"AL": "Albania",
		"DZ": "Algeria",
		"AS": "American Samoa",
		"AD": "Andorra",
		"AO": "Angola",
		"AI": "Anguilla",
		"AQ": "Antarctica",
		"AG": "Antigua and Barbuda",
		"AR": "Argentina",
		"AM": "Armenia",
		"AW": "Aruba",
		"AU": "Australia",
		"AT": "Austria",
		"AZ": "Azerbaijan",
		"BS": "Bahamas",
		"BH": "Bahrain",
		"BD": "Bangladesh",
		"BB": "Barbados",
		"BY": "Belarus",
		"BE": "Belgium",
		"BZ": "Belize",
		"BJ": "Benin",
		"BM": "Bermuda",
		"BT": "Bhutan",
		"BO": "Bolivia",
		"BA": "Bosnia and Herzegovina",
		"BW": "Botswana",
		"BV": "Bouvet Island",
		"BR": "Brazil",
		"IO": "British Indian Ocean Territory",
		"VG": "British Virgin Islands",
		"BN": "Brunei",
		"BG": "Bulgaria",
		"BF": "Burkina Faso",
		"BI": "Burundi",
		"KH": "Cambodia",
		"CM": "Cameroon",
		"CA": "Canada",
		"CV": "Cape Verde",
		"KY": "Cayman Islands",
		"CF": "Central African Republic",
		"TD": "Chad",
		"CL": "Chile",
		"CN": "China",
		"CX": "Christmas Island",
		"CC": "Cocos [Keeling] Islands",
		"CO": "Colombia",
		"KM": "Comoros",
		"CG": "Congo - Brazzaville",
		"CD": "Congo - Kinshasa",
		"CK": "Cook Islands",
		"CR": "Costa Rica",
		"CI": "Côte d’Ivoire",
		"HR": "Croatia",
		"CU": "Cuba",
		"CY": "Cyprus",
		"CZ": "Czech Republic",
		"DK": "Denmark",
		"DJ": "Djibouti",
		"DM": "Dominica",
		"DO": "Dominican Republic",
		"EC": "Ecuador",
		"EG": "Egypt",
		"SV": "El Salvador",
		"GQ": "Equatorial Guinea",
		"ER": "Eritrea",
		"EE": "Estonia",
		"ET": "Ethiopia",
		"FK": "Falkland Islands",
		"FO": "Faroe Islands",
		"FJ": "Fiji",
		"FI": "Finland",
		"FR": "France",
		"GF": "French Guiana",
		"PF": "French Polynesia",
		"TF": "French Southern Territories",
		"GA": "Gabon",
		"GM": "Gambia",
		"GE": "Georgia",
		"DE": "Germany",
		"GH": "Ghana",
		"GI": "Gibraltar",
		"GR": "Greece",
		"GL": "Greenland",
		"GD": "Grenada",
		"GP": "Guadeloupe",
		"GU": "Guam",
		"GT": "Guatemala",
		"GG": "Guernsey",
		"GN": "Guinea",
		"GW": "Guinea-Bissau",
		"GY": "Guyana",
		"HT": "Haiti",
		"HM": "Heard Island and McDonald Islands",
		"HN": "Honduras",
		"HK": "Hong Kong SAR China",
		"HU": "Hungary",
		"IS": "Iceland",
		"IN": "India",
		"ID": "Indonesia",
		"IR": "Iran",
		"IQ": "Iraq",
		"IE": "Ireland",
		"IM": "Isle of Man",
		"IL": "Israel",
		"IT": "Italy",
		"JM": "Jamaica",
		"JP": "Japan",
		"JE": "Jersey",
		"JO": "Jordan",
		"KZ": "Kazakhstan",
		"KE": "Kenya",
		"KI": "Kiribati",
		"KW": "Kuwait",
		"KG": "Kyrgyzstan",
		"LA": "Laos",
		"LV": "Latvia",
		"LB": "Lebanon",
		"LS": "Lesotho",
		"LR": "Liberia",
		"LY": "Libya",
		"LI": "Liechtenstein",
		"LT": "Lithuania",
		"LU": "Luxembourg",
		"MO": "Macau SAR China",
		"MK": "Macedonia",
		"MG": "Madagascar",
		"MW": "Malawi",
		"MY": "Malaysia",
		"MV": "Maldives",
		"ML": "Mali",
		"MT": "Malta",
		"MH": "Marshall Islands",
		"MQ": "Martinique",
		"MR": "Mauritania",
		"MU": "Mauritius",
		"YT": "Mayotte",
		"MX": "Mexico",
		"FM": "Micronesia",
		"MD": "Moldova",
		"MC": "Monaco",
		"MN": "Mongolia",
		"ME": "Montenegro",
		"MS": "Montserrat",
		"MA": "Morocco",
		"MZ": "Mozambique",
		"MM": "Myanmar [Burma]",
		"NA": "Namibia",
		"NR": "Nauru",
		"NP": "Nepal",
		"NL": "Netherlands",
		"AN": "Netherlands Antilles",
		"NC": "New Caledonia",
		"NZ": "New Zealand",
		"NI": "Nicaragua",
		"NE": "Niger",
		"NG": "Nigeria",
		"NU": "Niue",
		"NF": "Norfolk Island",
		"MP": "Northern Mariana Islands",
		"KP": "North Korea",
		"NO": "Norway",
		"OM": "Oman",
		"PK": "Pakistan",
		"PW": "Palau",
		"PS": "Palestinian Territories",
		"PA": "Panama",
		"PG": "Papua New Guinea",
		"PY": "Paraguay",
		"PE": "Peru",
		"PH": "Philippines",
		"PN": "Pitcairn Islands",
		"PL": "Poland",
		"PT": "Portugal",
		"PR": "Puerto Rico",
		"QA": "Qatar",
		"RE": "Réunion",
		"RO": "Romania",
		"RU": "Russia",
		"RW": "Rwanda",
		"BL": "Saint Barthélemy",
		"SH": "Saint Helena",
		"KN": "Saint Kitts and Nevis",
		"LC": "Saint Lucia",
		"MF": "Saint Martin",
		"PM": "Saint Pierre and Miquelon",
		"VC": "Saint Vincent and the Grenadines",
		"WS": "Samoa",
		"SM": "San Marino",
		"ST": "São Tomé and Príncipe",
		"SA": "Saudi Arabia",
		"SN": "Senegal",
		"RS": "Serbia",
		"SC": "Seychelles",
		"SL": "Sierra Leone",
		"SG": "Singapore",
		"SK": "Slovakia",
		"SI": "Slovenia",
		"SB": "Solomon Islands",
		"SO": "Somalia",
		"ZA": "South Africa",
		"GS": "South Georgia and the South Sandwich Islands",
		"KR": "South Korea",
		"ES": "Spain",
		"LK": "Sri Lanka",
		"SD": "Sudan",
		"SR": "Suriname",
		"SJ": "Svalbard and Jan Mayen",
		"SZ": "Swaziland",
		"SE": "Sweden",
		"CH": "Switzerland",
		"SY": "Syria",
		"TW": "Taiwan",
		"TJ": "Tajikistan",
		"TZ": "Tanzania",
		"TH": "Thailand",
		"TL": "Timor-Leste",
		"TG": "Togo",
		"TK": "Tokelau",
		"TO": "Tonga",
		"TT": "Trinidad and Tobago",
		"TN": "Tunisia",
		"TR": "Turkey",
		"TM": "Turkmenistan",
		"TC": "Turks and Caicos Islands",
		"TV": "Tuvalu",
		"UG": "Uganda",
		"UA": "Ukraine",
		"AE": "United Arab Emirates",
		"GB": "United Kingdom",
		"US": "United States",
		"UY": "Uruguay",
		"UM": "U.S. Minor Outlying Islands",
		"VI": "U.S. Virgin Islands",
		"UZ": "Uzbekistan",
		"VU": "Vanuatu",
		"VA": "Vatican City",
		"VE": "Venezuela",
		"VN": "Vietnam",
		"WF": "Wallis and Futuna",
		"EH": "Western Sahara",
		"YE": "Yemen",
		"ZM": "Zambia",
		"ZW": "Zimbabwe",
	}
)
