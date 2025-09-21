package constants

const (
	QueryInsertAnimal = `
		INSERT INTO animals 
			(id,farm_id,name,species,breed,birth_date,origin_type,origin_from,purchased_price,gender,weight_kg,health_status,registration_date,is_active,notes)
		VALUES
			($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)
		RETURNING id,farm_id,species,breed,birth_date,origin_type,origin_from,purchased_price,gender,weight_kg,health_status,registration_date,is_active,notes;
	`
	QueryDeleteAnimal = `
		DELETE FROM animals WHERE id = $1
	`
)
