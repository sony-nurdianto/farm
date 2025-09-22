package animalquery

const (
	QueryTotalAnimals = `
	SELECT COUNT(*) AS total
	FROM animals
	WHERE 
		farm_id = $1
		AND (
			COALESCE($2, '') = '' OR name ILIKE '%' || $2 || '%'
		) 
	`
	QueryGetAnimalASC = `
	SELECT 
		(id,farm_id,name,species,breed,birth_date,origin_type,origin_from,purchased_price,gender,weight_kg,health_status,registration_date,is_active,notes)
	FROM animals	
	WHERE 
		farm_id = $1
		AND (
			COALESCE($2, '') = '' OR name ILIKE '%' || $2 || '%'
		)
	ORDER BY registered_at ASC
	LIMIT $3
	OFFSET $4
	`
	QueryGetAnimalDESC = `
	SELECT 
		(id,farm_id,name,species,breed,birth_date,origin_type,origin_from,purchased_price,gender,weight_kg,health_status,registration_date,is_active,notes)
	FROM animals	
	WHERE 
		farm_id = $1
		AND (
			COALESCE($2, '') = '' OR name ILIKE '%' || $2 || '%'
		)
	ORDER BY registered_at DESC
	LIMIT $3
	OFFSET $4
	`
)
