#!/usr/bin/env bash

echo ""
echo "Running Hurl tests..."

# moves up to root from
if [[ -f "../.env" ]]; then
  # echo "Please run this script from the project root."
  cd ..
fi

if [[ ! -f ".env" ]]; then
  echo "Unable to find '.env' file."
  exit 1
fi

# source .env variables
source .env

# run user tests with admin token for testing
hurl \
  --variable craig_email=craig@gmail.com \
  --variable craig_password=@ssword472 \
  --secret super_admin_token=$SUPER_ADMIN_TOKEN \
  --test \
  test/users.hurl

# run admin tests with admin token for testing
hurl \
  --variable lisa_email=lisa@gmail.com \
  --variable lisa_password=Growl1ng! \
  --variable 1_species_name="Pilea peperomioides" \
  --variable 1_human_poison_toxic=false \
  --variable 1_pet_poison_toxic=false \
  --variable 1_human_edible=false \
  --variable 1_pet_edible=false \
  --variable 2_species_name="Crassula ovata" \
  --variable 2_human_poison_toxic=true \
  --variable 2_pet_poison_toxic=true \
  --variable 2_human_edible=false \
  --variable 2_pet_edible=false \
  --variable 1a_species_common_name="Chinese money plant" \
  --variable 1a_species_common_langcode="en" \
  --variable 1b_species_common_name="UFO plant" \
  --variable 1b_species_common_langcode="en" \
  --variable 1c_species_common_name="lefse plant" \
  --variable 1c_species_common_langcode="en" \
  --variable 1d_species_common_name="planta china del dinero" \
  --variable 1d_species_common_langcode="es" \
  --variable 1e_species_common_name="planta lefse" \
  --variable 1e_species_common_langcode="es" \
  --variable 1f_species_common_name="planta ONVI" \
  --variable 1f_species_common_langcode="es" \
  --variable 2a_species_common_name="jade plant" \
  --variable 2a_species_common_langcode="en" \
  --variable 2b_species_common_name="lucky plant" \
  --variable 2b_species_common_langcode="en" \
  --variable 2c_species_common_name="money plant" \
  --variable 2c_species_common_langcode="en" \
  --variable 2d_species_common_name="árbol de jade" \
  --variable 2d_species_common_langcode="es" \
  --variable 2e_species_common_name="Monedita" \
  --variable 2e_species_common_langcode="es" \
  --variable 2f_species_common_name="árbol de las monedas" \
  --variable 2f_species_common_langcode="es" \
  --variable 1_plant_type_name="Tropical" \
  --variable 1_plant_type_description="Tropical plants thrive in warm, humid environments and are often characterized by their lush, green foliage. They typically require consistent moisture and indirect light." \
  --variable 1_plant_type_maxtc=35 \
  --variable 1_plant_type_mintc=10 \
  --variable 1_plant_type_maxph=80 \
  --variable 1_plant_type_minph=30 \
  --variable 1_plant_type_soilom="2 parts Peat moss" \
  --variable 1_plant_type_soilgm="1 part Perlite" \
  --variable 1_plant_type_soildm="1 part Sand" \
  --variable 2_plant_type_name="Temperate" \
  --variable 2_plant_type_description="Temperate plants are adapted to regions with distinct seasons and moderate temperatures. They often have a dormant period and can tolerate cooler temperatures." \
  --variable 2_plant_type_maxtc=30 \
  --variable 2_plant_type_mintc=10 \
  --variable 2_plant_type_maxph=70 \
  --variable 2_plant_type_minph=30 \
  --variable 2_plant_type_soilom="2 parts Potting soil" \
  --variable 2_plant_type_soilgm="1 part Perlite" \
  --variable 2_plant_type_soildm="1 part Sand" \
  --variable 3_plant_type_name="Semi-Arid" \
  --variable 3_plant_type_description="Semi-arid plants are adapted to regions with moderate rainfall and can tolerate some drought. They often have succulent or waxy leaves to retain moisture." \
  --variable 3_plant_type_maxtc=30 \
  --variable 3_plant_type_mintc=10 \
  --variable 3_plant_type_maxph=60 \
  --variable 3_plant_type_minph=20 \
  --variable 3_plant_type_soilom="1 part Cactus mix" \
  --variable 3_plant_type_soilgm="1 part Pumice" \
  --variable 3_plant_type_soildm="1 part Sand" \
  --variable 4_plant_type_name="Arid" \
  --variable 4_plant_type_description="Arid plants are adapted to extremely dry environments and can tolerate prolonged periods without water. They often have thick, waxy leaves or spines to conserve moisture." \
  --variable 4_plant_type_maxtc=40 \
  --variable 4_plant_type_mintc=5 \
  --variable 4_plant_type_maxph=50 \
  --variable 4_plant_type_minph=10 \
  --variable 4_plant_type_soilom="1 part Cactus mix" \
  --variable 4_plant_type_soilgm="1 part Pumice" \
  --variable 4_plant_type_soildm="2 parts Sand" \
  --secret super_admin_token=$SUPER_ADMIN_TOKEN \
  --test \
  test/admin.hurl
