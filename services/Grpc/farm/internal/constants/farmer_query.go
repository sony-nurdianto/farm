package constants

const (
	QueryGetFarmAsc = `
		SELECT 
		    f.id,
		    f.farmer_id,
		    f.farm_name,
		    f.farm_type,
		    f.farm_size,
		    f.farm_status,
		    f.description,
		    f.created_at,
		    f.updated_at,
		    f.address_id,
		    a.street,
		    a.village,
		    a.sub_district,
		    a.city,
		    a.province,
		    a.postal_code
		FROM farms f
		LEFT JOIN addresses a ON f.address_id = a.id
		WHERE 
		    f.farmer_id = $1
		    AND (
		        COALESCE($2, '') = '' OR f.farm_name ILIKE '%' || $2 || '%'
		    )
		ORDER BY f.created_at ASC
		LIMIT $3 
		OFFSET $4
	`

	QueryGetFarmDesc = `
		SELECT 
		    f.id,
		    f.farmer_id,
		    f.farm_name,
		    f.farm_type,
		    f.farm_size,
		    f.farm_status,
		    f.description,
		    f.created_at,
		    f.updated_at,
		    f.address_id,
		    a.street,
		    a.village,
		    a.sub_district,
		    a.city,
		    a.province,
		    a.postal_code
		FROM farms f
		LEFT JOIN addresses a ON f.address_id = a.id
		WHERE 
		    f.farmer_id = $1
		    AND (
		        COALESCE($2, '') = '' OR f.farm_name ILIKE '%' || $2 || '%'
		    )
		ORDER BY f.created_at DESC
		LIMIT $3 
		OFFSET $4
	`

	QueryTotalFarm = `
	SELECT COUNT(*) as total
	FROM farms 
	WHERE farmer_id = $1 
  AND
		(
			COALESCE($2, '') = '' OR farm_name ILIKE '%' || $2 || '%'
		)
	`

	QueryGetFarmByID = `
		SELECT 
		    f.id,
		    f.farmer_id,
		    f.farm_name,
		    f.farm_type,
		    f.farm_size,
		    f.farm_status,
		    f.description,
		    f.created_at,
		    f.updated_at,
		    f.address_id,
		    a.street,
		    a.village,
		    a.sub_district,
		    a.city,
		    a.province,
		    a.postal_code
		FROM farms f
		LEFT JOIN addresses a ON f.address_id = a.id
		WHERE f.id = $1
	`
)
