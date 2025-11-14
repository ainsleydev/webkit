interface Language {
	code: string;
	name: string;
	encoding: string;
}

export const languages: Language[] = [
	{
		code: 'as_IN',
		name: 'Assamese (India)',
		encoding: 'UTF-8',
	},
	{
		code: 'fy_NL',
		name: 'Western Frisian (Netherlands)',
		encoding: 'UTF-8',
	},
	{
		code: 'gl_ES',
		name: 'Galician (Spain)',
		encoding: 'ISO-8859-15',
	},
	{
		code: 'he_IL',
		name: 'Hebrew (Israel)',
		encoding: 'ISO-8859-8',
	},
	{
		code: 'kk_KZ',
		name: 'Kazakh (Kazakhstan)',
		encoding: 'PT154',
	},
	{
		code: 'nl_NL',
		name: 'Dutch (Netherlands)',
		encoding: 'ISO-8859-15',
	},
	{
		code: 'ti_ER',
		name: 'Tigrinya (Eritrea)',
		encoding: 'UTF-8',
	},
	{
		code: 'cs',
		name: 'Czech',
		encoding: 'ISO-8859-2',
	},
	{
		code: 'gu_IN',
		name: 'Gujarati (India)',
		encoding: 'UTF-8',
	},
	{
		code: 'nso_ZA',
		name: 'Pedi (South Africa)',
		encoding: 'UTF-8',
	},
	{
		code: 'ki',
		name: 'Kikuyu',
		encoding: 'UTF-8',
	},
	{
		code: 'ki_KE',
		name: 'Kikuyu (Kenya)',
		encoding: 'UTF-8',
	},
	{
		code: 'vi',
		name: 'Vietnamese',
		encoding: 'UTF-8',
	},
	{
		code: 'ks_IN',
		name: 'Kashmiri (India)',
		encoding: 'UTF-8',
	},
	{
		code: 'ca_AD',
		name: 'Catalan (Andorra)',
		encoding: 'ISO-8859-15',
	},
	{
		code: 'en_NG',
		name: 'UTF-8',
		encoding: 'UTF-8',
	},
	{
		code: 'es_PY',
		name: 'Spanish (Paraguay)',
		encoding: 'ISO-8859-1',
	},
	{
		code: 'fr_FR',
		name: 'French (France)',
		encoding: 'ISO-8859-15',
	},
	{
		code: 'sl_SI',
		name: 'Slovenian (Slovenia)',
		encoding: 'ISO-8859-2',
	},
	{
		code: 'ta_IN',
		name: 'Tamil (India)',
		encoding: 'UTF-8',
	},
	{
		code: 'bm_ML',
		name: 'Bambara (Mali)',
		encoding: 'UTF-8',
	},
	{
		code: 'eu_ES',
		name: 'Basque (Spain)',
		encoding: 'ISO-8859-15',
	},
	{
		code: 'pa',
		name: 'Punjabi',
		encoding: 'UTF-8',
	},
	{
		code: 'uz_Arab',
		name: 'Uzbek (Arabic)',
		encoding: 'UTF-8',
	},
	{
		code: 'zh_Hans',
		name: 'Chinese (Simplified Han)',
		encoding: 'GB2312',
	},
	{
		code: 'ha_Latn_GH',
		name: 'Hausa (Latin, Ghana)',
		encoding: 'UTF-8',
	},
	{
		code: 'khq_ML',
		name: 'Koyra Chiini (Mali)',
		encoding: 'UTF-8',
	},
	{
		code: 'pap_AN',
		name: 'Papiamento (Netherlands)',
		encoding: 'UTF-8',
	},
	{
		code: 'sid_ET',
		name: 'Sidamo (Ethiopia)',
		encoding: 'UTF-8',
	},
	{
		code: 'zh_Hant_TW',
		name: 'Chinese (Traditional Han, Taiwan)',
		encoding: 'GB2312',
	},
	{
		code: 'eo',
		name: 'Esperanto',
		encoding: 'UTF-8',
	},
	{
		code: 'ha_Latn_NE',
		name: 'Hausa (Latin, Niger)',
		encoding: 'UTF-8',
	},
	{
		code: 'hi',
		name: 'Hindi',
		encoding: 'UTF-8',
	},
	{
		code: 'sl',
		name: 'Slovenian',
		encoding: 'ISO-8859-2',
	},
	{
		code: 'sr_Cyrl',
		name: 'Serbian (Cyrillic)',
		encoding: 'UTF-8',
	},
	{
		code: 'ug_CN',
		name: 'Uighur (China)',
		encoding: 'UTF-8',
	},
	{
		code: 'fr_CG',
		name: 'French (Congo - Brazzaville)',
		encoding: 'ISO-8859-15',
	},
	{
		code: 'mhr_RU',
		name: 'Eastern Mari (Russia)',
		encoding: 'UTF-8',
	},
	{
		code: 'so_KE',
		name: 'Somali (Kenya)',
		encoding: 'ISO-8859-1',
	},
	{
		code: 'te',
		name: 'Telugu',
		encoding: 'UTF-8',
	},
	{
		code: 'vun',
		name: 'Vunjo',
		encoding: 'UTF-8',
	},
	{
		code: 'ak_GH',
		name: 'Akan (Ghana)',
		encoding: 'UTF-8',
	},
	{
		code: 'fr_RE',
		name: 'French (Réunion)',
		encoding: 'ISO-8859-15',
	},
	{
		code: 'st_ZA',
		name: 'Southern Sotho (South Africa)',
		encoding: 'ISO-8859-1',
	},
	{
		code: 'en_IE',
		name: 'English (Ireland)',
		encoding: 'ISO-8859-15',
	},
	{
		code: 'om_KE',
		name: 'Oromo (Kenya)',
		encoding: 'ISO-8859-1',
	},
	{
		code: 'es_SV',
		name: 'Spanish (El Salvador)',
		encoding: 'ISO-8859-1',
	},
	{
		code: 'luo',
		name: 'Luo',
		encoding: 'UTF-8',
	},
	{
		code: 'lv',
		name: 'Latvian',
		encoding: 'ISO-8859-13',
	},
	{
		code: 'ru_UA',
		name: 'Russian (Ukraine)',
		encoding: 'KOI8-U',
	},
	{
		code: 'ses',
		name: 'Koyraboro Senni',
		encoding: 'UTF-8',
	},
	{
		code: 'shi_Tfng',
		name: 'Tachelhit (Tifinagh)',
		encoding: 'UTF-8',
	},
	{
		code: 'cy',
		name: 'Welsh',
		encoding: 'ISO-8859-14',
	},
	{
		code: 'it_CH',
		name: 'Italian (Switzerland)',
		encoding: 'ISO-8859-1',
	},
	{
		code: 'saq',
		name: 'Samburu',
		encoding: 'UTF-8',
	},
	{
		code: 'fr_CM',
		name: 'French (Cameroon)',
		encoding: 'ISO-8859-15',
	},
	{
		code: 'ha_NG',
		name: 'Hausa (Nigeria)',
		encoding: 'UTF-8',
	},
	{
		code: 'or',
		name: 'Oriya',
		encoding: 'UTF-8',
	},
	{
		code: 'sa_IN',
		name: 'Sanskrit (India)',
		encoding: 'UTF-8',
	},
	{
		code: 'da_DK',
		name: 'Danish (Denmark)',
		encoding: 'ISO-8859-1',
	},
	{
		code: 'pl',
		name: 'Polish',
		encoding: 'ISO-8859-2',
	},
	{
		code: 'sw',
		name: 'Swahili',
		encoding: 'UTF-8',
	},
	{
		code: 'pa_Arab_PK',
		name: 'Punjabi (Arabic, Pakistan)',
		encoding: 'UTF-8',
	},
	{
		code: 'aa_ER',
		name: 'Afar (Eritrea)',
		encoding: 'UTF-8',
	},
	{
		code: 'ar_IQ',
		name: 'Arabic (Iraq)',
		encoding: 'ISO-8859-6',
	},
	{
		code: 'en_DK',
		name: 'English (Denmark)',
		encoding: 'ISO-8859-1',
	},
	{
		code: 'fr_ML',
		name: 'French (Mali)',
		encoding: 'ISO-8859-15',
	},
	{
		code: 'mr',
		name: 'Marathi',
		encoding: 'UTF-8',
	},
	{
		code: 'nds_NL',
		name: 'Low German (Netherlands)',
		encoding: 'UTF-8',
	},
	{
		code: 'or_IN',
		name: 'Oriya (India)',
		encoding: 'UTF-8',
	},
	{
		code: 'shi_Tfng_MA',
		name: 'Tachelhit (Tifinagh, Morocco)',
		encoding: 'UTF-8',
	},
	{
		code: 'to',
		name: 'Tonga',
		encoding: 'UTF-8',
	},
	{
		code: 'lv_LV',
		name: 'Latvian (Latvia)',
		encoding: 'ISO-8859-13',
	},
	{
		code: 'zh_HK',
		name: 'Chinese (Hong Kong)',
		encoding: 'BIG5-HKSCS',
	},
	{
		code: 'az_Cyrl',
		name: 'Azerbaijani (Cyrillic)',
		encoding: 'UTF-8',
	},
	{
		code: 'nd',
		name: 'North Ndebele',
		encoding: 'UTF-8',
	},
	{
		code: 'pt_PT',
		name: 'Portuguese (Portugal)',
		encoding: 'ISO-8859-15',
	},
	{
		code: 'ee',
		name: 'Ewe',
		encoding: 'UTF-8',
	},
	{
		code: 'fr_MG',
		name: 'French (Madagascar)',
		encoding: 'ISO-8859-15',
	},
	{
		code: 'fr_SN',
		name: 'French (Senegal)',
		encoding: 'ISO-8859-15',
	},
	{
		code: 'nb',
		name: 'Norwegian Bokmål',
		encoding: 'ISO-8859-1',
	},
	{
		code: 'ne_NP',
		name: 'Nepali (Nepal)',
		encoding: 'UTF-8',
	},
	{
		code: 'ar',
		name: 'Arabic',
		encoding: 'ISO-8859-6',
	},
	{
		code: 'en_AG',
		name: 'English (Antigua \u0026 Barbuda)',
		encoding: 'UTF-8',
	},
	{
		code: 'fr_NE',
		name: 'French (Niger)',
		encoding: 'ISO-8859-15',
	},
	{
		code: 'ja_JP',
		name: 'Japanese (Japan)',
		encoding: 'UTF-8',
	},
	{
		code: 'km',
		name: 'Khmer',
		encoding: 'UTF-8',
	},
	{
		code: 'tzm_Latn',
		name: 'Central Morocco Tamazight (Latin)',
		encoding: 'UTF-8',
	},
	{
		code: 'zu_ZA',
		name: 'Zulu (South Africa)',
		encoding: 'ISO-8859-1',
	},
	{
		code: 'fil_PH',
		name: 'Filipino (Philippines)',
		encoding: 'UTF-8',
	},
	{
		code: 'sv_SE',
		name: 'Swedish (Sweden)',
		encoding: 'ISO-8859-1',
	},
	{
		code: 'en_SG',
		name: 'English (Singapore)',
		encoding: 'ISO-8859-1',
	},
	{
		code: 'sr_Cyrl_ME',
		name: 'Serbian (Cyrillic, Montenegro)',
		encoding: 'UTF-8',
	},
	{
		code: 'fr_BE',
		name: 'French (Belgium)',
		encoding: 'ISO-8859-15',
	},
	{
		code: 'fr_CI',
		name: 'French (Côte d’Ivoire)',
		encoding: 'ISO-8859-15',
	},
	{
		code: 'id_ID',
		name: 'Indonesian (Indonesia)',
		encoding: 'ISO-8859-1',
	},
	{
		code: 'sq_AL',
		name: 'Albanian (Albania)',
		encoding: 'ISO-8859-1',
	},
	{
		code: 'dv_MV',
		name: 'Dhivehi (Maldives)',
		encoding: 'UTF-8',
	},
	{
		code: 'rw',
		name: 'Kinyarwanda',
		encoding: 'UTF-8',
	},
	{
		code: 'sd_IN',
		name: 'Sindhi (India)',
		encoding: 'UTF-8',
	},
	{
		code: 'yo_NG',
		name: 'Yoruba (Nigeria)',
		encoding: 'UTF-8',
	},
	{
		code: 'af',
		name: 'Afrikaans',
		encoding: 'UTF-8',
	},
	{
		code: 'bg_BG',
		name: 'Bulgarian (Bulgaria)',
		encoding: 'CP1251',
	},
	{
		code: 'en_VI',
		name: 'English (U.S. Virgin Islands)',
		encoding: 'UTF-8',
	},
	{
		code: 'es_MX',
		name: 'Spanish (Mexico)',
		encoding: 'ISO-8859-1',
	},
	{
		code: 'fil',
		name: 'Filipino',
		encoding: 'UTF-8',
	},
	{
		code: 'es_UY',
		name: 'Spanish (Uruguay)',
		encoding: 'ISO-8859-1',
	},
	{
		code: 'gv_GB',
		name: 'Manx (United Kingdom)',
		encoding: 'ISO-8859-1',
	},
	{
		code: 'ar_SD',
		name: 'Arabic (Sudan)',
		encoding: 'ISO-8859-6',
	},
	{
		code: 'fr_TG',
		name: 'French (Togo)',
		encoding: 'ISO-8859-15',
	},
	{
		code: 'kk_Cyrl_KZ',
		name: 'Kazakh (Cyrillic, Kazakhstan)',
		encoding: 'PT154',
	},
	{
		code: 'ps_AF',
		name: 'Pashto (Afghanistan)',
		encoding: 'UTF-8',
	},
	{
		code: 'shi',
		name: 'Tachelhit',
		encoding: 'UTF-8',
	},
	{
		code: 'sn',
		name: 'Shona',
		encoding: 'UTF-8',
	},
	{
		code: 'es_ES',
		name: 'Spanish (Spain)',
		encoding: 'ISO-8859-15',
	},
	{
		code: 'sg',
		name: 'Sango',
		encoding: 'UTF-8',
	},
	{
		code: 'tt_RU',
		name: 'Tatar (Russia)',
		encoding: 'UTF-8',
	},
	{
		code: 'zu',
		name: 'Zulu',
		encoding: 'ISO-8859-1',
	},
	{
		code: 'hu',
		name: 'Hungarian',
		encoding: 'ISO-8859-2',
	},
	{
		code: 'hu_HU',
		name: 'Hungarian (Hungary)',
		encoding: 'ISO-8859-2',
	},
	{
		code: 'ii',
		name: 'Sichuan Yi',
		encoding: 'UTF-8',
	},
	{
		code: 'lg_UG',
		name: 'Ganda (Uganda)',
		encoding: 'ISO-8859-10',
	},
	{
		code: 'aa_ET',
		name: 'Afar (Ethiopia)',
		encoding: 'UTF-8',
	},
	{
		code: 'nl',
		name: 'Dutch',
		encoding: 'UTF-8',
	},
	{
		code: 'rof',
		name: 'Rombo',
		encoding: 'UTF-8',
	},
	{
		code: 'ar_OM',
		name: 'Arabic (Oman)',
		encoding: 'ISO-8859-6',
	},
	{
		code: 'en_GB',
		name: 'English (United Kingdom)',
		encoding: 'ISO-8859-1',
	},
	{
		code: 'es_EC',
		name: 'Spanish (Ecuador)',
		encoding: 'ISO-8859-1',
	},
	{
		code: 'ml_IN',
		name: 'Malayalam (India)',
		encoding: 'UTF-8',
	},
	{
		code: 'sk',
		name: 'Slovak',
		encoding: 'ISO-8859-2',
	},
	{
		code: 'tr_TR',
		name: 'Turkish (Turkey)',
		encoding: 'ISO-8859-9',
	},
	{
		code: 'kl_GL',
		name: 'Kalaallisut (Greenland)',
		encoding: 'ISO-8859-1',
	},
	{
		code: 'mai_IN',
		name: 'Maithili (India)',
		encoding: 'UTF-8',
	},
	{
		code: 'ur_IN',
		name: 'Urdu (India)',
		encoding: 'UTF-8',
	},
	{
		code: 'de_LI',
		name: 'German (Liechtenstein)',
		encoding: 'ISO-8859-15',
	},
	{
		code: 'ha_Latn',
		name: 'Hausa (Latin)',
		encoding: 'UTF-8',
	},
	{
		code: 'hsb_DE',
		name: 'Upper Sorbian (Germany)',
		encoding: 'UTF-8',
	},
	{
		code: 'ne',
		name: 'Nepali',
		encoding: 'UTF-8',
	},
	{
		code: 'pt_GW',
		name: 'Portuguese (Guinea-Bissau)',
		encoding: 'ISO-8859-15',
	},
	{
		code: 'sq',
		name: 'Albanian',
		encoding: 'ISO-8859-1',
	},
	{
		code: 'tr_CY',
		name: 'Turkish (Cyprus)',
		encoding: 'ISO-8859-9',
	},
	{
		code: 'byn_ER',
		name: 'Bilin (Eritrea)',
		encoding: 'UTF-8',
	},
	{
		code: 'ebu',
		name: 'Embu',
		encoding: 'UTF-8',
	},
	{
		code: 'pa_Guru_IN',
		name: 'Punjabi (Gurmukhi, India)',
		encoding: 'UTF-8',
	},
	{
		code: 'ebu_KE',
		name: 'Embu (Kenya)',
		encoding: 'UTF-8',
	},
	{
		code: 'he',
		name: 'Hebrew',
		encoding: 'ISO-8859-8',
	},
	{
		code: 'ss_ZA',
		name: 'Swati (South Africa)',
		encoding: 'UTF-8',
	},
	{
		code: 'uz_Latn_UZ',
		name: 'Uzbek (Latin, Uzbekistan)',
		encoding: 'UTF-8',
	},
	{
		code: 'mas_KE',
		name: 'Masai (Kenya)',
		encoding: 'UTF-8',
	},
	{
		code: 'mk',
		name: 'Macedonian',
		encoding: 'ISO-8859-5',
	},
	{
		code: 'de_DE',
		name: 'German (Germany)',
		encoding: 'ISO-8859-15',
	},
	{
		code: 'ha_Latn_NG',
		name: 'Hausa (Latin, Nigeria)',
		encoding: 'UTF-8',
	},
	{
		code: 'bho_IN',
		name: 'Bhojpuri (India)',
		encoding: 'UTF-8',
	},
	{
		code: 'es_US',
		name: 'Spanish (United States)',
		encoding: 'ISO-8859-1',
	},
	{
		code: 'kam_KE',
		name: 'Kamba (Kenya)',
		encoding: 'UTF-8',
	},
	{
		code: 'tr',
		name: 'Turkish',
		encoding: 'ISO-8859-9',
	},
	{
		code: 'es_NI',
		name: 'Spanish (Nicaragua)',
		encoding: 'ISO-8859-1',
	},
	{
		code: 'kea',
		name: 'Kabuverdianu',
		encoding: 'UTF-8',
	},
	{
		code: 'ur_PK',
		name: 'Urdu (Pakistan)',
		encoding: 'UTF-8',
	},
	{
		code: 'bem_ZM',
		name: 'Bemba (Zambia)',
		encoding: 'UTF-8',
	},
	{
		code: 'chr',
		name: 'Cherokee',
		encoding: 'UTF-8',
	},
	{
		code: 'sr_Latn_RS',
		name: 'Serbian (Latin, Serbia)',
		encoding: 'UTF-8',
	},
	{
		code: 'sr_ME',
		name: 'Serbian (Montenegro)',
		encoding: 'UTF-8',
	},
	{
		code: 'ca_FR',
		name: 'Catalan (France)',
		encoding: 'ISO-8859-15',
	},
	{
		code: 'en_ZA',
		name: 'English (South Africa)',
		encoding: 'ISO-8859-1',
	},
	{
		code: 'es_AR',
		name: 'Spanish (Argentina)',
		encoding: 'ISO-8859-1',
	},
	{
		code: 'fr_BL',
		name: 'French (Saint Barthélemy)',
		encoding: 'ISO-8859-15',
	},
	{
		code: 'pa_PK',
		name: 'Panjabi (Pakistan)',
		encoding: 'UTF-8',
	},
	{
		code: 'pt',
		name: 'Portuguese',
		encoding: 'ISO-8859-1',
	},
	{
		code: 'fa',
		name: 'Persian',
		encoding: 'UTF-8',
	},
	{
		code: 'is',
		name: 'Icelandic',
		encoding: 'ISO-8859-1',
	},
	{
		code: 'rw_RW',
		name: 'Kinyarwanda (Rwanda)',
		encoding: 'UTF-8',
	},
	{
		code: 'ti_ET',
		name: 'Tigrinya (Ethiopia)',
		encoding: 'UTF-8',
	},
	{
		code: 'guz',
		name: 'Gusii',
		encoding: 'UTF-8',
	},
	{
		code: 'mfe_MU',
		name: 'Morisyen (Mauritius)',
		encoding: 'UTF-8',
	},
	{
		code: 'nn',
		name: 'Norwegian Nynorsk',
		encoding: 'ISO-8859-1',
	},
	{
		code: 'teo_KE',
		name: 'Teso (Kenya)',
		encoding: 'UTF-8',
	},
	{
		code: 'zh_Hans_MO',
		name: 'Chinese (Simplified Han, Macau SAR China)',
		encoding: 'GB2312',
	},
	{
		code: 'af_ZA',
		name: 'Afrikaans (South Africa)',
		encoding: 'ISO-8859-1',
	},
	{
		code: 'ee_TG',
		name: 'Ewe (Togo)',
		encoding: 'UTF-8',
	},
	{
		code: 'fr_MQ',
		name: 'French (Martinique)',
		encoding: 'ISO-8859-15',
	},
	{
		code: 'lt_LT',
		name: 'Lithuanian (Lithuania)',
		encoding: 'ISO-8859-13',
	},
	{
		code: 'pa_Guru',
		name: 'Punjabi (Gurmukhi)',
		encoding: 'UTF-8',
	},
	{
		code: 'saq_KE',
		name: 'Samburu (Kenya)',
		encoding: 'UTF-8',
	},
	{
		code: 'xh_ZA',
		name: 'Xhosa South Africa)',
		encoding: 'ISO-8859-1',
	},
	{
		code: 'zh_TW',
		name: 'Chinese (Taiwan)',
		encoding: 'BIG5',
	},
	{
		code: 'bn',
		name: 'Bengali',
		encoding: 'UTF-8',
	},
	{
		code: 'mas',
		name: 'Masai',
		encoding: 'UTF-8',
	},
	{
		code: 'ms',
		name: 'Malay',
		encoding: 'ISO-8859-1',
	},
	{
		code: 'my_MM',
		name: 'Burmese (Myanmar [Burma])',
		encoding: 'UTF-8',
	},
	{
		code: 'shi_Latn_MA',
		name: 'Tachelhit (Latin, Morocco)',
		encoding: 'UTF-8',
	},
	{
		code: 'ta_LK',
		name: 'Tamil (Sri Lanka)',
		encoding: 'UTF-8',
	},
	{
		code: 'ti',
		name: 'Tigrinya',
		encoding: 'UTF-8',
	},
	{
		code: 'ik_CA',
		name: 'Inupiaq (Canada)',
		encoding: 'UTF-8',
	},
	{
		code: 'lag',
		name: 'Langi',
		encoding: 'UTF-8',
	},
	{
		code: 'tzm',
		name: 'Central Morocco Tamazight',
		encoding: 'UTF-8',
	},
	{
		code: 'yi_US',
		name: 'Yiddish (America)',
		encoding: 'CP1255',
	},
	{
		code: 'fr_BF',
		name: 'French (Burkina Faso)',
		encoding: 'ISO-8859-15',
	},
	{
		code: 'teo',
		name: 'Teso',
		encoding: 'UTF-8',
	},
	{
		code: 'ar_BH',
		name: 'Arabic (Bahrain)',
		encoding: 'ISO-8859-6',
	},
	{
		code: 'en_NA',
		name: 'English (Namibia)',
		encoding: 'UTF-8',
	},
	{
		code: 'hy',
		name: 'Armenian',
		encoding: 'ARMSCII-8',
	},
	{
		code: 'lo_LA',
		name: 'Lao (Laos)',
		encoding: 'UTF-8',
	},
	{
		code: 'sr_RS',
		name: 'Serbian (Serbia)',
		encoding: 'UTF-8',
	},
	{
		code: 'sv_FI',
		name: 'Swedish (Finland)',
		encoding: 'ISO-8859-15',
	},
	{
		code: 'bem',
		name: 'Bemba',
		encoding: 'UTF-8',
	},
	{
		code: 'es_GQ',
		name: 'Spanish (Equatorial Guinea)',
		encoding: 'ISO-8859-1',
	},
	{
		code: 'gl',
		name: 'Galician',
		encoding: 'ISO-8859-15',
	},
	{
		code: 'sr_Latn_BA',
		name: 'Serbian (Latin, Bosnia and Herzegovina)',
		encoding: 'UTF-8',
	},
	{
		code: 'ts_ZA',
		name: 'Tsonga (South Africa)',
		encoding: 'UTF-8',
	},
	{
		code: 'seh_MZ',
		name: 'Sena (Mozambique)',
		encoding: 'UTF-8',
	},
	{
		code: 'az_AZ',
		name: 'Azerbaijani (Azerbaijan)',
		encoding: 'UTF-8',
	},
	{
		code: 'en_AU',
		name: 'English (Australia)',
		encoding: 'ISO-8859-1',
	},
	{
		code: 'fa_IR',
		name: 'Persian (Iran)',
		encoding: 'UTF-8',
	},
	{
		code: 'gu',
		name: 'Gujarati',
		encoding: 'UTF-8',
	},
	{
		code: 'haw_US',
		name: 'Hawaiian (United States)',
		encoding: 'UTF-8',
	},
	{
		code: 'ml',
		name: 'Malayalam',
		encoding: 'UTF-8',
	},
	{
		code: 'nan_TW',
		name: 'Min Nan Chinese (Taiwan)',
		encoding: 'UTF-8',
	},
	{
		code: 'zh_SG',
		name: 'Chinese (Singapore)',
		encoding: 'GB2312',
	},
	{
		code: 'bez_TZ',
		name: 'Bena (Tanzania)',
		encoding: 'UTF-8',
	},
	{
		code: 'ka_GE',
		name: 'Georgian (Georgia)',
		encoding: 'GEORGIAN-PS',
	},
	{
		code: 'kw',
		name: 'Cornish',
		encoding: 'ISO-8859-1',
	},
	{
		code: 'aa_DJ',
		name: 'Afar (Ddjibouti)',
		encoding: 'ISO-8859-1',
	},
	{
		code: 'ar_QA',
		name: 'Arabic (Qatar)',
		encoding: 'ISO-8859-6',
	},
	{
		code: 'br_FR',
		name: 'Breton (France)',
		encoding: 'ISO-8859-15',
	},
	{
		code: 'luy',
		name: 'Luyia',
		encoding: 'UTF-8',
	},
	{
		code: 'tig_ER',
		name: 'Tigre (Eritrea)',
		encoding: 'UTF-8',
	},
	{
		code: 'hne_IN',
		name: 'Chhattisgarhi (India)',
		encoding: 'UTF-8',
	},
	{
		code: 'mer_KE',
		name: 'Meru (Kenya)',
		encoding: 'UTF-8',
	},
	{
		code: 'nd_ZW',
		name: 'North Ndebele (Zimbabwe)',
		encoding: 'UTF-8',
	},
	{
		code: 'ne_IN',
		name: 'Nepali (India)',
		encoding: 'UTF-8',
	},
	{
		code: 'pa_Arab',
		name: 'Punjabi (Arabic)',
		encoding: 'UTF-8',
	},
	{
		code: 'te_IN',
		name: 'Telugu (India)',
		encoding: 'UTF-8',
	},
	{
		code: 'ar_MA',
		name: 'Arabic (Morocco)',
		encoding: 'ISO-8859-6',
	},
	{
		code: 'es',
		name: 'Spanish',
		encoding: 'ISO-8859-1',
	},
	{
		code: 'ko',
		name: 'Korean',
		encoding: 'UTF-8',
	},
	{
		code: 'xog',
		name: 'Soga',
		encoding: 'UTF-8',
	},
	{
		code: 'shi_Latn',
		name: 'Tachelhit (Latin)',
		encoding: 'UTF-8',
	},
	{
		code: 'so_ET',
		name: 'Somali (Ethiopia)',
		encoding: 'UTF-8',
	},
	{
		code: 'sv',
		name: 'Swedish',
		encoding: 'ISO-8859-15',
	},
	{
		code: 'fr_BJ',
		name: 'French (Benin)',
		encoding: 'ISO-8859-15',
	},
	{
		code: 'seh',
		name: 'Sena',
		encoding: 'UTF-8',
	},
	{
		code: 'bez',
		name: 'Bena',
		encoding: 'UTF-8',
	},
	{
		code: 'fr_TD',
		name: 'French (Chad)',
		encoding: 'ISO-8859-15',
	},
	{
		code: 'mk_MK',
		name: 'Macedonian (Macedonia)',
		encoding: 'ISO-8859-5',
	},
	{
		code: 'mt',
		name: 'Maltese',
		encoding: 'ISO-8859-3',
	},
	{
		code: 'sr_Cyrl_RS',
		name: 'Serbian (Cyrillic, Serbia)',
		encoding: 'UTF-8',
	},
	{
		code: 'tzm_Latn_MA',
		name: 'Central Morocco Tamazight (Latin, Morocco)',
		encoding: 'UTF-8',
	},
	{
		code: 'ar_KW',
		name: 'Arabic (Kuwait)',
		encoding: 'ISO-8859-6',
	},
	{
		code: 'en_AS',
		name: 'English (American Samoa)',
		encoding: 'UTF-8',
	},
	{
		code: 'en_BZ',
		name: 'English (Belize)',
		encoding: 'UTF-8',
	},
	{
		code: 'es_VE',
		name: 'Spanish (Venezuela)',
		encoding: 'ISO-8859-1',
	},
	{
		code: 'fr_GP',
		name: 'French (Guadeloupe)',
		encoding: 'ISO-8859-15',
	},
	{
		code: 'om',
		name: 'Oromo',
		encoding: 'UTF-8',
	},
	{
		code: 'so_DJ',
		name: 'Somali (Djibouti)',
		encoding: 'ISO-8859-1',
	},
	{
		code: 'fur_IT',
		name: 'Friulian (Italy)',
		encoding: 'UTF-8',
	},
	{
		code: 'haw',
		name: 'Hawaiian',
		encoding: 'UTF-8',
	},
	{
		code: 'mr_IN',
		name: 'Marathi (India)',
		encoding: 'UTF-8',
	},
	{
		code: 'rwk_TZ',
		name: 'Rwa (Tanzania)',
		encoding: 'UTF-8',
	},
	{
		code: 'az',
		name: 'Azerbaijani',
		encoding: 'UTF-8',
	},
	{
		code: 'kab',
		name: 'Kabyle',
		encoding: 'UTF-8',
	},
	{
		code: 'to_TO',
		name: 'Tonga (Tonga)',
		encoding: 'UTF-8',
	},
	{
		code: 'bn_IN',
		name: 'Bengali (India)',
		encoding: 'UTF-8',
	},
	{
		code: 'de_LU',
		name: 'German (Luxembourg)',
		encoding: 'ISO-8859-15',
	},
	{
		code: 'dz_BT',
		name: 'Dzongkha (Buthan)',
		encoding: 'UTF-8',
	},
	{
		code: 'fr_CH',
		name: 'French (Switzerland)',
		encoding: 'ISO-8859-1',
	},
	{
		code: 'ar_LB',
		name: 'Arabic (Lebanon)',
		encoding: 'ISO-8859-6',
	},
	{
		code: 'bs_BA',
		name: 'Bosnian (Bosnia and Herzegovina)',
		encoding: 'ISO-8859-2',
	},
	{
		code: 'ca',
		name: 'Catalan',
		encoding: 'ISO-8859-15',
	},
	{
		code: 'es_PA',
		name: 'Spanish (Panama)',
		encoding: 'ISO-8859-1',
	},
	{
		code: 'am',
		name: 'Amharic',
		encoding: 'UTF-8',
	},
	{
		code: 'es_DO',
		name: 'Spanish (Dominican Republic)',
		encoding: 'ISO-8859-1',
	},
	{
		code: 'hr_HR',
		name: 'Croatian (Croatia)',
		encoding: 'ISO-8859-2',
	},
	{
		code: 'kde',
		name: 'Makonde',
		encoding: 'UTF-8',
	},
	{
		code: 'ig_NG',
		name: 'Igbo (Nigeria)',
		encoding: 'UTF-8',
	},
	{
		code: 'kln',
		name: 'Kalenjin',
		encoding: 'UTF-8',
	},
	{
		code: 'ku_TR',
		name: 'Kurdish (Turkey)',
		encoding: 'ISO-8859-9',
	},
	{
		code: 'pa_IN',
		name: 'Panjabi (India)',
		encoding: 'UTF-8',
	},
	{
		code: 'fi_FI',
		name: 'Finnish (Finland)',
		encoding: 'ISO-8859-15',
	},
	{
		code: 'gsw',
		name: 'Swiss German',
		encoding: 'UTF-8',
	},
	{
		code: 'khq',
		name: 'Koyra Chiini',
		encoding: 'UTF-8',
	},
	{
		code: 'ru_MD',
		name: 'Russian (Moldova)',
		encoding: 'ISO-8859-5',
	},
	{
		code: 'shs_CA',
		name: 'Shuswap (Canada)',
		encoding: 'UTF-8',
	},
	{
		code: 'uz_Latn',
		name: 'Uzbek (Latin)',
		encoding: 'UTF-8',
	},
	{
		code: 'yo',
		name: 'Yoruba',
		encoding: 'UTF-8',
	},
	{
		code: 'es_PE',
		name: 'Spanish (Peru)',
		encoding: 'ISO-8859-1',
	},
	{
		code: 'fo',
		name: 'Faroese',
		encoding: 'ISO-8859-1',
	},
	{
		code: 'ro',
		name: 'Romanian',
		encoding: 'ISO-8859-2',
	},
	{
		code: 'az_Latn_AZ',
		name: 'Azerbaijani (Latin, Azerbaijan)',
		encoding: 'UTF-8',
	},
	{
		code: 'csb_PL',
		name: 'Kashubian (Poland)',
		encoding: 'UTF-8',
	},
	{
		code: 'fr_RW',
		name: 'French (Rwanda)',
		encoding: 'ISO-8859-15',
	},
	{
		code: 'nl_AW',
		name: 'Dutch (Aruba)',
		encoding: 'UTF-8',
	},
	{
		code: 'ro_RO',
		name: 'Romanian (Romania)',
		encoding: 'ISO-8859-2',
	},
	{
		code: 'se_NO',
		name: 'Northern Sami (Norway)',
		encoding: 'UTF-8',
	},
	{
		code: 'sk_SK',
		name: 'Slovak (Slovakia)',
		encoding: 'ISO-8859-2',
	},
	{
		code: 'zh_Hant',
		name: 'Chinese (Traditional Han)',
		encoding: 'GB2312',
	},
	{
		code: 'it_IT',
		name: 'Italian (Italy)',
		encoding: 'ISO-8859-15',
	},
	{
		code: 'lb_LU',
		name: 'Luxembourgish (Luxembourg)',
		encoding: 'UTF-8',
	},
	{
		code: 'mn_MN',
		name: 'Mongolian (Mongolia)',
		encoding: 'UTF-8',
	},
	{
		code: 'tg_TJ',
		name: 'Tajik (Tajikistan)',
		encoding: 'KOI8-T',
	},
	{
		code: 'wae_CH',
		name: 'Walser (Switzerland)',
		encoding: 'UTF-8',
	},
	{
		code: 'ar_DZ',
		name: 'Arabic (Algeria)',
		encoding: 'ISO-8859-6',
	},
	{
		code: 'id',
		name: 'Indonesian',
		encoding: 'ISO-8859-1',
	},
	{
		code: 'ht_HT',
		name: 'Haitian (Haiti)',
		encoding: 'UTF-8',
	},
	{
		code: 'pt_MZ',
		name: 'Portuguese (Mozambique)',
		encoding: 'ISO-8859-15',
	},
	{
		code: 'ru',
		name: 'Russian',
		encoding: 'ISO-8859-5',
	},
	{
		code: 'en',
		name: 'English',
		encoding: 'UTF-8',
	},
	{
		code: 'ro_MD',
		name: 'Romanian (Moldova)',
		encoding: 'ISO-8859-2',
	},
	{
		code: 'cy_GB',
		name: 'Welsh (United Kingdom)',
		encoding: 'ISO-8859-14',
	},
	{
		code: 'gez_ET',
		name: 'Geez (Ethiopia)',
		encoding: 'UTF-8',
	},
	{
		code: 'mag_IN',
		name: 'Magahi (India)',
		encoding: 'UTF-8',
	},
	{
		code: 'nb_NO',
		name: 'Norwegian Bokmål (Norway)',
		encoding: 'ISO-8859-1',
	},
	{
		code: 'sr_Latn_ME',
		name: 'Serbian (Latin, Montenegro)',
		encoding: 'UTF-8',
	},
	{
		code: 'zh_Hant_MO',
		name: 'Chinese (Traditional Han, Macau SAR China)',
		encoding: 'GB2312',
	},
	{
		code: 'en_HK',
		name: 'English (Hong Kong SAR China)',
		encoding: 'ISO-8859-1',
	},
	{
		code: 'en_NZ',
		name: 'English (New Zealand)',
		encoding: 'ISO-8859-1',
	},
	{
		code: 'fr_MC',
		name: 'French (Monaco)',
		encoding: 'ISO-8859-15',
	},
	{
		code: 'kl',
		name: 'Kalaallisut',
		encoding: 'ISO-8859-1',
	},
	{
		code: 'luo_KE',
		name: 'Luo (Kenya)',
		encoding: 'UTF-8',
	},
	{
		code: 'nds_DE',
		name: 'Low German (German)',
		encoding: 'UTF-8',
	},
	{
		code: 'sq_MK',
		name: 'Albanian (Macedonia)',
		encoding: 'UTF-8',
	},
	{
		code: 'ar_YE',
		name: 'Arabic (Yemen)',
		encoding: 'ISO-8859-6',
	},
	{
		code: 'el',
		name: 'Greek',
		encoding: 'ISO-8859-7',
	},
	{
		code: 'rwk',
		name: 'Rwa',
		encoding: 'UTF-8',
	},
	{
		code: 'zh',
		name: 'Chinese',
		encoding: 'GB2312',
	},
	{
		code: 'ky_KG',
		name: 'Kirghiz (Kyrgyzstan)',
		encoding: 'UTF-8',
	},
	{
		code: 'uk',
		name: 'Ukrainian',
		encoding: 'KOI8-U',
	},
	{
		code: 'xog_UG',
		name: 'Soga (Uganda)',
		encoding: 'UTF-8',
	},
	{
		code: 'be',
		name: 'Belarusian',
		encoding: 'UTF-8',
	},
	{
		code: 'bo_IN',
		name: 'Tibetan (India)',
		encoding: 'UTF-8',
	},
	{
		code: 'cgg_UG',
		name: 'Chiga (Uganda)',
		encoding: 'UTF-8',
	},
	{
		code: 'de_BE',
		name: 'German (Belgium)',
		encoding: 'ISO-8859-15',
	},
	{
		code: 'fa_AF',
		name: 'Persian (Afghanistan)',
		encoding: 'UTF-8',
	},
	{
		code: 'ga_IE',
		name: 'Irish (Ireland)',
		encoding: 'ISO-8859-15',
	},
	{
		code: 'th_TH',
		name: 'Thai (Thailand)',
		encoding: 'TIS-620',
	},
	{
		code: 'tl_PH',
		name: 'Tagalog (Philippines)',
		encoding: 'ISO-8859-1',
	},
	{
		code: 'en_MT',
		name: 'English (Malta)',
		encoding: 'UTF-8',
	},
	{
		code: 'fr_CF',
		name: 'French (Central African Republic)',
		encoding: 'ISO-8859-15',
	},
	{
		code: 'fr_LU',
		name: 'French (Luxembourg)',
		encoding: 'ISO-8859-15',
	},
	{
		code: 'kln_KE',
		name: 'Kalenjin (Kenya)',
		encoding: 'UTF-8',
	},
	{
		code: 'kw_GB',
		name: 'Cornish (United Kingdom)',
		encoding: 'ISO-8859-1',
	},
	{
		code: 'mas_TZ',
		name: 'Masai (Tanzania)',
		encoding: 'UTF-8',
	},
	{
		code: 'ta',
		name: 'Tamil',
		encoding: 'UTF-8',
	},
	{
		code: 'kok',
		name: 'Konkani',
		encoding: 'UTF-8',
	},
	{
		code: 'mer',
		name: 'Meru',
		encoding: 'UTF-8',
	},
	{
		code: 'ar_TN',
		name: 'Arabic (Tunisia)',
		encoding: 'ISO-8859-6',
	},
	{
		code: 'ee_GH',
		name: 'Ewe (Ghana)',
		encoding: 'UTF-8',
	},
	{
		code: 'fr_DJ',
		name: 'French (Djibouti)',
		encoding: 'ISO-8859-15',
	},
	{
		code: 'fr_GA',
		name: 'French (Gabon)',
		encoding: 'ISO-8859-15',
	},
	{
		code: 'kde_TZ',
		name: 'Makonde (Tanzania)',
		encoding: 'UTF-8',
	},
	{
		code: 'lt',
		name: 'Lithuanian',
		encoding: 'ISO-8859-13',
	},
	{
		code: 'dav_KE',
		name: 'Taita (Kenya)',
		encoding: 'UTF-8',
	},
	{
		code: 'fr_CA',
		name: 'French (Canada)',
		encoding: 'ISO-8859-15',
	},
	{
		code: 'ga',
		name: 'Irish',
		encoding: 'ISO-8859-15',
	},
	{
		code: 'ko_KR',
		name: 'Korean (South Korea)',
		encoding: 'UTF-8',
	},
	{
		code: 'nyn',
		name: 'Nyankole',
		encoding: 'UTF-8',
	},
	{
		code: 'vi_VN',
		name: 'Vietnamese (Vietnam)',
		encoding: 'UTF-8',
	},
	{
		code: 'zh_CN',
		name: 'Chinese (China)',
		encoding: 'GB2312',
	},
	{
		code: 'en_BE',
		name: 'English (Belgium)',
		encoding: 'UTF-8',
	},
	{
		code: 'es_419',
		name: 'Spanish (Latin America)',
		encoding: 'ISO-8859-1',
	},
	{
		code: 'es_CL',
		name: 'Spanish (Chile)',
		encoding: 'ISO-8859-1',
	},
	{
		code: 'ha',
		name: 'Hausa',
		encoding: 'UTF-8',
	},
	{
		code: 've_ZA',
		name: 'Venda (South Africa)',
		encoding: 'UTF-8',
	},
	{
		code: 'cv_RU',
		name: 'Chuvash (Russia)',
		encoding: 'UTF-8',
	},
	{
		code: 'en_PK',
		name: 'English (Pakistan)',
		encoding: 'UTF-8',
	},
	{
		code: 'ii_CN',
		name: 'Sichuan Yi (China)',
		encoding: 'UTF-8',
	},
	{
		code: 'wal_ET',
		name: 'Wolaytta (Ethiopia)',
		encoding: 'UTF-8',
	},
	{
		code: 'gd_GB',
		name: 'Scottish Gaelic (United Kingdom)',
		encoding: 'ISO-8859-15',
	},
	{
		code: 'si_LK',
		name: 'Sinhala (Sri Lanka)',
		encoding: 'UTF-8',
	},
	{
		code: 'de',
		name: 'German',
		encoding: 'ISO-8859-15',
	},
	{
		code: 'en_ZW',
		name: 'English (Zimbabwe)',
		encoding: 'ISO-8859-1',
	},
	{
		code: 'ar_AE',
		name: 'Arabic (United Arab Emirates)',
		encoding: 'ISO-8859-6',
	},
	{
		code: 'ff',
		name: 'Fulah',
		encoding: 'UTF-8',
	},
	{
		code: 'ff_SN',
		name: 'Fulah (Senegal)',
		encoding: 'UTF-8',
	},
	{
		code: 'sw_TZ',
		name: 'Swahili (Tanzania)',
		encoding: 'UTF-8',
	},
	{
		code: 'unm_US',
		name: 'Unami (America)',
		encoding: 'UTF-8',
	},
	{
		code: 'ms_MY',
		name: 'Malay (Malaysia)',
		encoding: 'ISO-8859-1',
	},
	{
		code: 'so_SO',
		name: 'Somali (Somalia)',
		encoding: 'ISO-8859-1',
	},
	{
		code: 'sr',
		name: 'Serbian',
		encoding: 'UTF-8',
	},
	{
		code: 'da',
		name: 'Danish',
		encoding: 'ISO-8859-1',
	},
	{
		code: 'el_CY',
		name: 'Greek (Cyprus)',
		encoding: 'ISO-8859-7',
	},
	{
		code: 'en_MP',
		name: 'English (Northern Mariana Islands)',
		encoding: 'UTF-8',
	},
	{
		code: 'es_GT',
		name: 'Spanish (Guatemala)',
		encoding: 'ISO-8859-1',
	},
	{
		code: 'lag_TZ',
		name: 'Langi (Tanzania)',
		encoding: 'UTF-8',
	},
	{
		code: 'so',
		name: 'Somali',
		encoding: 'ISO-8859-1',
	},
	{
		code: 'uk_UA',
		name: 'Ukrainian (Ukraine)',
		encoding: 'KOI8-U',
	},
	{
		code: 'ar_SY',
		name: 'Arabic (Syria)',
		encoding: 'ISO-8859-6',
	},
	{
		code: 'es_BO',
		name: 'Spanish (Bolivia)',
		encoding: 'ISO-8859-1',
	},
	{
		code: 'fy_DE',
		name: 'Western Frisian (Germany)',
		encoding: 'UTF-8',
	},
	{
		code: 'li_NL',
		name: 'Limburgan (Netherlands)',
		encoding: 'UTF-8',
	},
	{
		code: 'mi_NZ',
		name: 'Maori (New Zealand)',
		encoding: 'ISO-8859-13',
	},
	{
		code: 'ak',
		name: 'Akan',
		encoding: 'UTF-8',
	},
	{
		code: 'bg',
		name: 'Bulgarian',
		encoding: 'CP1251',
	},
	{
		code: 'en_MH',
		name: 'English (Marshall Islands)',
		encoding: 'UTF-8',
	},
	{
		code: 'fr_GN',
		name: 'French (Guinea)',
		encoding: 'ISO-8859-15',
	},
	{
		code: 'lg',
		name: 'Ganda',
		encoding: 'ISO-8859-10',
	},
	{
		code: 'naq_NA',
		name: 'Nama (Namibia)',
		encoding: 'UTF-8',
	},
	{
		code: 'om_ET',
		name: 'Oromo (Ethiopia)',
		encoding: 'UTF-8',
	},
	{
		code: 'ar_SA',
		name: 'Arabic (Saudi Arabia)',
		encoding: 'ISO-8859-6',
	},
	{
		code: 'my',
		name: 'Burmese',
		encoding: 'UTF-8',
	},
	{
		code: 'ar_LY',
		name: 'Arabic (Libya)',
		encoding: 'ISO-8859-6',
	},
	{
		code: 'ber_DZ',
		name: 'UTF-8',
		encoding: 'UTF-8',
	},
	{
		code: 'ca_ES',
		name: 'Catalan (Spain)',
		encoding: 'ISO-8859-15',
	},
	{
		code: 'fi',
		name: 'Finnish',
		encoding: 'ISO-8859-15',
	},
	{
		code: 'guz_KE',
		name: 'Gusii (Kenya)',
		encoding: 'UTF-8',
	},
	{
		code: 'ja',
		name: 'Japanese',
		encoding: 'UTF-8',
	},
	{
		code: 'jmc',
		name: 'Machame',
		encoding: 'UTF-8',
	},
	{
		code: 'ig',
		name: 'Igbo',
		encoding: 'UTF-8',
	},
	{
		code: 'es_HN',
		name: 'Spanish (Honduras)',
		encoding: 'ISO-8859-1',
	},
	{
		code: 'et_EE',
		name: 'Estonian (Estonia)',
		encoding: 'ISO-8859-15',
	},
	{
		code: 'fr',
		name: 'French',
		encoding: 'ISO-8859-15',
	},
	{
		code: 'li_BE',
		name: 'Limburgan (Belgium)',
		encoding: 'UTF-8',
	},
	{
		code: 'mfe',
		name: 'Morisyen',
		encoding: 'UTF-8',
	},
	{
		code: 'sg_CF',
		name: 'Sango (Central African Republic)',
		encoding: 'UTF-8',
	},
	{
		code: 'am_ET',
		name: 'Amharic (Ethiopia)',
		encoding: 'UTF-8',
	},
	{
		code: 'brx_IN',
		name: 'Bodo (India)',
		encoding: 'UTF-8',
	},
	{
		code: 'en_US',
		name: 'English (United States)',
		encoding: 'ISO-8859-1',
	},
	{
		code: 'gez_ER',
		name: 'Geez (Eritrea)',
		encoding: 'UTF-8',
	},
	{
		code: 'uz_Arab_AF',
		name: 'Uzbek (Arabic, Afghanistan)',
		encoding: 'UTF-8',
	},
	{
		code: 'ar_IN',
		name: 'Arabic (India)',
		encoding: 'UTF-8',
	},
	{
		code: 'as',
		name: 'Assamese',
		encoding: 'UTF-8',
	},
	{
		code: 'bs',
		name: 'Bosnian',
		encoding: 'ISO-8859-2',
	},
	{
		code: 'fr_GQ',
		name: 'French (Equatorial Guinea)',
		encoding: 'ISO-8859-15',
	},
	{
		code: 'gv',
		name: 'Manx',
		encoding: 'ISO-8859-1',
	},
	{
		code: 'vun_TZ',
		name: 'Vunjo (Tanzania)',
		encoding: 'UTF-8',
	},
	{
		code: 'wa_BE',
		name: 'Walloon (Belgium)',
		encoding: 'UTF-8',
	},
	{
		code: 'fr_BI',
		name: 'French (Burundi)',
		encoding: 'ISO-8859-15',
	},
	{
		code: 'nr_ZA',
		name: 'South Ndebele (South Africa)',
		encoding: 'UTF-8',
	},
	{
		code: 'sc_IT',
		name: 'Sardinian (Italy)',
		encoding: 'UTF-8',
	},
	{
		code: 'zh_Hans_SG',
		name: 'Chinese (Simplified Han, Singapore)',
		encoding: 'GB2312',
	},
	{
		code: 'af_NA',
		name: 'Afrikaans (Namibia)',
		encoding: 'UTF-8',
	},
	{
		code: 'oc_FR',
		name: 'Occitan (France)',
		encoding: 'ISO-8859-1',
	},
	{
		code: 'crh_UA',
		name: 'Crimean Tatar (Ukraine)',
		encoding: 'UTF-8',
	},
	{
		code: 'kk',
		name: 'Kazakh',
		encoding: 'PT154',
	},
	{
		code: 'kok_IN',
		name: 'Konkani (India)',
		encoding: 'UTF-8',
	},
	{
		code: 'th',
		name: 'Thai',
		encoding: 'TIS-620',
	},
	{
		code: 'uz_UZ',
		name: 'Uzbek (Uzbekistan)',
		encoding: 'UTF-8',
	},
	{
		code: 'az_Cyrl_AZ',
		name: 'Azerbaijani (Cyrillic, Azerbaijan)',
		encoding: 'UTF-8',
	},
	{
		code: 'is_IS',
		name: 'Icelandic (Iceland)',
		encoding: 'ISO-8859-1',
	},
	{
		code: 'kn',
		name: 'Kannada',
		encoding: 'UTF-8',
	},
	{
		code: 'rof_TZ',
		name: 'Rombo (Tanzania)',
		encoding: 'UTF-8',
	},
	{
		code: 'ms_BN',
		name: 'Malay (Brunei)',
		encoding: 'ISO-8859-1',
	},
	{
		code: 'dav',
		name: 'Taita',
		encoding: 'UTF-8',
	},
	{
		code: 'en_ZM',
		name: 'English (Zambia)',
		encoding: 'UTF-8',
	},
	{
		code: 'si',
		name: 'Sinhala',
		encoding: 'UTF-8',
	},
	{
		code: 'tn_ZA',
		name: 'Tswana (South Africa)',
		encoding: 'UTF-8',
	},
	{
		code: 'ka',
		name: 'Georgian',
		encoding: 'GEORGIAN-PS',
	},
	{
		code: 'ar_JO',
		name: 'Arabic (Jordan)',
		encoding: 'ISO-8859-6',
	},
	{
		code: 'es_CO',
		name: 'Spanish (Colombia)',
		encoding: 'ISO-8859-1',
	},
	{
		code: 'kea_CV',
		name: 'Kabuverdianu (Cape Verde)',
		encoding: 'UTF-8',
	},
	{
		code: 'lij_IT',
		name: 'Ligurian (Italy)',
		encoding: 'UTF-8',
	},
	{
		code: 'ps',
		name: 'Pashto',
		encoding: 'UTF-8',
	},
	{
		code: 'ber_MA',
		name: 'Berber (Morocco)',
		encoding: 'UTF-8',
	},
	{
		code: 'fr_MF',
		name: 'French (Saint Martin)',
		encoding: 'ISO-8859-15',
	},
	{
		code: 'hr',
		name: 'Croatian',
		encoding: 'ISO-8859-2',
	},
	{
		code: 'ru_RU',
		name: 'Russian (Russia)',
		encoding: 'ISO-8859-5',
	},
	{
		code: 'fo_FO',
		name: 'Faroese (Faroe Islands)',
		encoding: 'ISO-8859-1',
	},
	{
		code: 'iw_IL',
		name: 'Hebrew (Israel)',
		encoding: 'ISO-8859-8',
	},
	{
		code: 'jmc_TZ',
		name: 'Machame (Tanzania)',
		encoding: 'UTF-8',
	},
	{
		code: 'kn_IN',
		name: 'Kannada (India)',
		encoding: 'UTF-8',
	},
	{
		code: 'mg',
		name: 'Malagasy',
		encoding: 'ISO-8859-15',
	},
	{
		code: 'naq',
		name: 'Nama',
		encoding: 'UTF-8',
	},
	{
		code: 'en_CA',
		name: 'English (Canada)',
		encoding: 'ISO-8859-1',
	},
	{
		code: 'fr_CD',
		name: 'French (Congo - Kinshasa)',
		encoding: 'ISO-8859-15',
	},
	{
		code: 'ur',
		name: 'Urdu',
		encoding: 'UTF-8',
	},
	{
		code: 'wo_SN',
		name: 'Wolof (Senegal)',
		encoding: 'UTF-8',
	},
	{
		code: 'az_Latn',
		name: 'Azerbaijani (Latin)',
		encoding: 'UTF-8',
	},
	{
		code: 'bn_BD',
		name: 'Bengali (Bangladesh)',
		encoding: 'UTF-8',
	},
	{
		code: 'de_AT',
		name: 'German (Austria)',
		encoding: 'ISO-8859-15',
	},
	{
		code: 'nn_NO',
		name: 'Norwegian Nynorsk (Norway)',
		encoding: 'ISO-8859-1',
	},
	{
		code: 'yue_HK',
		name: 'Yue Chinese (China)',
		encoding: 'UTF-8',
	},
	{
		code: 'en_TT',
		name: 'English (Trinidad and Tobago)',
		encoding: 'UTF-8',
	},
	{
		code: 'es_CR',
		name: 'Spanish (Costa Rica)',
		encoding: 'ISO-8859-1',
	},
	{
		code: 'pl_PL',
		name: 'Polish (Poland)',
		encoding: 'ISO-8859-2',
	},
	{
		code: 'ast_ES',
		name: 'Asturian (Spain)',
		encoding: 'ISO-8859-15',
	},
	{
		code: 'bm',
		name: 'Bambara',
		encoding: 'UTF-8',
	},
	{
		code: 'cgg',
		name: 'Chiga',
		encoding: 'UTF-8',
	},
	{
		code: 'eu',
		name: 'Basque',
		encoding: 'ISO-8859-15',
	},
	{
		code: 'gsw_CH',
		name: 'Swiss German (Switzerland)',
		encoding: 'UTF-8',
	},
	{
		code: 'kk_Cyrl',
		name: 'Kazakh (Cyrillic)',
		encoding: 'PT154',
	},
	{
		code: 'nl_BE',
		name: 'Dutch (Belgium)',
		encoding: 'ISO-8859-15',
	},
	{
		code: 'rm_CH',
		name: 'Romansh (Switzerland)',
		encoding: 'UTF-8',
	},
	{
		code: 'sr_Latn',
		name: 'Serbian (Latin)',
		encoding: 'UTF-8',
	},
	{
		code: 'asa',
		name: 'Asu',
		encoding: 'UTF-8',
	},
	{
		code: 'chr_US',
		name: 'Cherokee (United States)',
		encoding: 'UTF-8',
	},
	{
		code: 'el_GR',
		name: 'Greek (Greece)',
		encoding: 'ISO-8859-7',
	},
	{
		code: 'nyn_UG',
		name: 'Nyankole (Uganda)',
		encoding: 'UTF-8',
	},
	{
		code: 'et',
		name: 'Estonian',
		encoding: 'ISO-8859-15',
	},
	{
		code: 'kab_DZ',
		name: 'Kabyle (Algeria)',
		encoding: 'UTF-8',
	},
	{
		code: 'rm',
		name: 'Romansh',
		encoding: 'UTF-8',
	},
	{
		code: 'asa_TZ',
		name: 'Asu (Tanzania)',
		encoding: 'UTF-8',
	},
	{
		code: 'be_BY',
		name: 'Belarusian (Belarus)',
		encoding: 'UTF-8',
	},
	{
		code: 'bo',
		name: 'Tibetan',
		encoding: 'UTF-8',
	},
	{
		code: 'de_CH',
		name: 'German (Switzerland)',
		encoding: 'ISO-8859-1',
	},
	{
		code: 'en_UM',
		name: 'English (U.S. Minor Outlying Islands)',
		encoding: 'UTF-8',
	},
	{
		code: 'ses_ML',
		name: 'Koyraboro Senni (Mali)',
		encoding: 'UTF-8',
	},
	{
		code: 'sw_KE',
		name: 'Swahili (Kenya)',
		encoding: 'UTF-8',
	},
	{
		code: 'en_MU',
		name: 'English (Mauritius)',
		encoding: 'UTF-8',
	},
	{
		code: 'iu_CA',
		name: 'Inuktitut (Canada)',
		encoding: 'UTF-8',
	},
	{
		code: 'os_RU',
		name: 'Ossetian (Russia)',
		encoding: 'UTF-8',
	},
	{
		code: 'zh_Hans_HK',
		name: 'Chinese (Simplified Han, Hong Kong SAR China)',
		encoding: 'BIG5-HKSCS',
	},
	{
		code: 'an_ES',
		name: 'Aragonese (Spain)',
		encoding: 'ISO-8859-15',
	},
	{
		code: 'bo_CN',
		name: 'Tibetan (China)',
		encoding: 'UTF-8',
	},
	{
		code: 'tk_TM',
		name: 'Turkmen (Turkmenistan)',
		encoding: 'UTF-8',
	},
	{
		code: 'cs_CZ',
		name: 'Czech (Czech Republic)',
		encoding: 'ISO-8859-2',
	},
	{
		code: 'hy_AM',
		name: 'Armenian (Armenia)',
		encoding: 'ARMSCII-8',
	},
	{
		code: 'sn_ZW',
		name: 'Shona (Zimbabwe)',
		encoding: 'UTF-8',
	},
	{
		code: 'ar_EG',
		name: 'Arabic (Egypt)',
		encoding: 'ISO-8859-6',
	},
	{
		code: 'mg_MG',
		name: 'Malagasy (Madagascar)',
		encoding: 'ISO-8859-15',
	},
	{
		code: 'teo_UG',
		name: 'Teso (Uganda)',
		encoding: 'UTF-8',
	},
	{
		code: 'uz',
		name: 'Uzbek',
		encoding: 'UTF-8',
	},
	{
		code: 'uz_Cyrl',
		name: 'Uzbek (Cyrillic)',
		encoding: 'UTF-8',
	},
	{
		code: 'en_PH',
		name: 'English (Philippines)',
		encoding: 'ISO-8859-1',
	},
	{
		code: 'hi_IN',
		name: 'Hindi (India)',
		encoding: 'UTF-8',
	},
	{
		code: 'it',
		name: 'Italian',
		encoding: 'ISO-8859-1',
	},
	{
		code: 'pt_BR',
		name: 'Portuguese (Brazil)',
		encoding: 'ISO-8859-1',
	},
	{
		code: 'zh_Hans_CN',
		name: 'Chinese (Simplified Han, China)',
		encoding: 'GB2312',
	},
	{
		code: 'en_GU',
		name: 'English (Guam)',
		encoding: 'UTF-8',
	},
	{
		code: 'en_IN',
		name: 'English (India)',
		encoding: 'UTF-8',
	},
	{
		code: 'en_JM',
		name: 'English (Jamaica)',
		encoding: 'UTF-8',
	},
	{
		code: 'es_CU',
		name: 'Spanish (Cuba)',
		encoding: 'UTF-8',
	},
	{
		code: 'kam',
		name: 'Kamba',
		encoding: 'UTF-8',
	},
	{
		code: 'sr_Cyrl_BA',
		name: 'Serbian (Cyrillic, Bosnia and Herzegovina)',
		encoding: 'UTF-8',
	},
	{
		code: 'uz_Cyrl_UZ',
		name: 'Uzbek (Cyrillic, Uzbekistan)',
		encoding: 'UTF-8',
	},
	{
		code: 'zh_Hant_HK',
		name: 'Chinese (Traditional Han, Hong Kong SAR China)',
		encoding: 'BIG5-HKSCS',
	},
	{
		code: 'ca_IT',
		name: 'Catalan (Italy)',
		encoding: 'ISO-8859-15',
	},
	{
		code: 'en_BW',
		name: 'English (Botswana)',
		encoding: 'ISO-8859-1',
	},
	{
		code: 'es_PR',
		name: 'Spanish (Puerto Rico)',
		encoding: 'ISO-8859-1',
	},
	{
		code: 'fr_KM',
		name: 'French (Comoros)',
		encoding: 'ISO-8859-15',
	},
	{
		code: 'km_KH',
		name: 'Khmer (Cambodia)',
		encoding: 'UTF-8',
	},
	{
		code: 'luy_KE',
		name: 'Luyia (Kenya)',
		encoding: 'UTF-8',
	},
	{
		code: 'mt_MT',
		name: 'Maltese (Malta)',
		encoding: 'ISO-8859-3',
	},
];
