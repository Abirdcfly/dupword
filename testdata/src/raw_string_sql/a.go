package raw_string_sql

func createQuery() string {
	const q = `
INSERT INTO addresses (
  contact_id,
  address_type_id,
  line1,
  line2,
  city,
  region,
  postal_code,
  country,
  is_deleted,
  version,
  created_at,
  updated_at,
  deleted_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, 0, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL);
`
	return q
}
