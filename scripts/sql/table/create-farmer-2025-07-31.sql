CREATE TABLE farmers (
    id UUID NOT NULL,
    full_name VARCHAR(225) NOT NULL,
    email VARCHAR(225) UNIQUE NOT NULL,
    phone VARCHAR(225) NOT NULL,
    registered_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    last_login TIMESTAMP WITH TIME ZONE,
    verified BOOLEAN DEFAULT FALSE,
    profile_photo_url TEXT,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
)
