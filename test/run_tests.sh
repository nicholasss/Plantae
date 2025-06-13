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
  --secret super_admin_token=$SUPER_ADMIN_TOKEN \
  --test \
  test/admin.hurl

# run user tests with admin token for testing
hurl \
  --variable craig_email=craig@gmail.com \
  --variable craig_password=@ssword472 \
  --secret super_admin_token=$SUPER_ADMIN_TOKEN \
  --test \
  test/users.hurl
