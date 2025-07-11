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
  --variable craig_lang_code=en \
  --secret super_admin_token=$SUPER_ADMIN_TOKEN \
  --jobs 1 \
  --test \
  test/users.hurl

# run admin tests with admin token for testing
hurl \
  --variable lisa_email=lisa@gmail.com \
  --variable lisa_password=Growl1ng! \
  --variable lisa_lang_code=en \
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
  --variable 1_light_name="Bright direct" \
  --variable 1_light_description="Unfiltered sunlight hitting the plant directly, like right in a south-facing window." \
  --variable 1_light_description_alt="Unfiltered sun exposure" \
  --variable 2_light_name="Bright indirect" \
  --variable 2_light_description="Bright light, but filtered or diffused, such as near a south-facing window with a sheer curtain, or a few feet away from an unobstructed east or west window." \
  --variable 2_light_description_alt="Bright, diffused light." \
  --variable 3_light_name="Medium indirect" \
  --variable 3_light_description="Good, consistent light, but never direct sun. Think a few feet back from an east or west window, or near a north-facing window." \
  --variable 3_light_description_alt="Consistent, non-direct light" \
  --variable 4_light_name="Low indirect" \
  --variable 4_light_description="Very little light, often found in the corner of a room, far from any window, or in a north-facing room with small windows." \
  --variable 4_light_description_alt="Very little natural light." \
  --variable 1_plant_water_type="Temperate" \
  --variable 1_plant_water_description="This plant is a temperate species that prefers its soil to dry out between watering sessions to avoid root rot. It thrives in typical indoor conditions." \
  --variable 1_plant_water_mm=50 \
  --variable 2_plant_water_type="Semi-Arid" \
  --variable 2_plant_water_description="As an arid succulent, the Jade plant stores water in its leaves and requires very infrequent watering. It is crucial to allow the soil to fully dry between waterings." \
  --variable 2_plant_water_days=15 \
  --secret super_admin_token=$SUPER_ADMIN_TOKEN \
  --jobs 1 \
  --test \
  test/admin.hurl

# run user tests for creating users plants
hurl \
  --variable craig_email=craig@gmail.com \
  --variable craig_password=@ssword472 \
  --variable craig_lang_code=en \
  --variable lisa_email=lisa@gmail.com \
  --variable lisa_password=Growl1ng! \
  --variable lisa_lang_code=en \
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
  --variable 1_plant_name="pepperoni" \
  --variable 1_plant_new_name="round guy" \
  --variable 1_plant_adoption="2020-01-02T00:00:00-05:00" \
  --variable 1_plant_new_adoption="2021-01-02T00:00:00-05:00" \
  --variable 2_plant_name="bonsai" \
  --variable 2_plant_new_name="mini tree" \
  --variable 2_plant_adoption="2020-11-30T00:00:00-05:00" \
  --variable 2_plant_new_adoption="2021-11-30T00:00:00-05:00" \
  --secret super_admin_token=$SUPER_ADMIN_TOKEN \
  --jobs 1 \
  --test \
  test/users_plants.hurl
