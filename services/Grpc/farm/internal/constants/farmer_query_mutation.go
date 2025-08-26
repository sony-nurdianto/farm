package constants

const (
	QueryInsertFarm = `
		insert into farms 
			(
				id,
				farmer_id, 
				farm_name,
				farm_type,
				farm_size,
				farm_status,
				description,
				address_id,
				created_at,
				updated_at	
			)
		values
			($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
		returning 
				id, 
				farmer_id, 
				farm_name, 
				farm_type, 
				farm_size,
				farm_status,
				description,
				address_id,
				created_at,
				updated_at; 
	`
	QueryInsertFarmAddress = `
		insert into addresses	
			(
				id,
				street,
				village,
				sub_district,
				city,
				province,
				postal_code,
				created_at,
				updated_at
			)
		values
			($1,$2,$3,$4,$5,$6,$7,$8,$9)
		returning 
			id,
			street,
			village,
			sub_district,
			city,
			province,
			postal_code,
			created_at,
			updated_at
	`

	QueryUpdateFarm = `
		update farm
		set 
				farm_name = coalesce($1,farm_name),
				farm_type = coalesce($2,farm_type),
				farm_size = coalesce($3,farm_size),
				farm_status = coalesce($4,farm_status),
				description = coalesce($5,description),
				updated_at = $6
		where id = $7
		returning 
				id, 
				farmer_id, 
				farm_name, 
				farm_type, 
				farm_size,
				farm_status,
				description,
				address_id,
				created_at,
				updated_at;
	`

	QueryUpdateFarmAddress = `
		update addresses
		set	
			street = coalesce($1,street),
			village = coalesce($2,village),
			sub_district = coalesce($3,sub_district),
			city = coalesce($4,city),
			province = coalesce($5,province),
			postal_code = coalesce($6,postal_code)
			updated_at = $7 
		where id = $8
		returning 
			id,
			street,
			village,
			sub_district,
			city,
			province,
			postal_code,
			created_at,
			updated_at
	`

	QueryDeleteFarm = `
		delete from farm where id = $1 cascade
	`
)
