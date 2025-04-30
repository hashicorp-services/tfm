#!/bin/bash


export SRC_TFE_HOSTNAME="select-newt.benjamin-lykins.sbx.hashidemos.io"
export SRC_TFE_TOKEN=""
export DST_TFC_HOSTNAME="ace-blowfish.benjamin-lykins.sbx.hashidemos.io"
export DST_TFC_ORG="default"
export DST_TFC_TOKEN=""

export SRC_ORGS="default,tfm-prereq-0,tfm-prereq-1,tfm-prereq-10,tfm-prereq-11,tfm-prereq-12,tfm-prereq-13,tfm-prereq-14,tfm-prereq-15,tfm-prereq-16,tfm-prereq-17,tfm-prereq-18,tfm-prereq-19,tfm-prereq-2,tfm-prereq-20,tfm-prereq-21,tfm-prereq-22,tfm-prereq-23,tfm-prereq-24,tfm-prereq-25,tfm-prereq-26,tfm-prereq-27,tfm-prereq-28,tfm-prereq-29,tfm-prereq-3,tfm-prereq-30,tfm-prereq-31,tfm-prereq-32,tfm-prereq-33,tfm-prereq-34,tfm-prereq-35,tfm-prereq-36,tfm-prereq-37,tfm-prereq-38,tfm-prereq-39,tfm-prereq-4,tfm-prereq-40,tfm-prereq-41,tfm-prereq-42,tfm-prereq-43,tfm-prereq-44,tfm-prereq-45,tfm-prereq-46,tfm-prereq-47,tfm-prereq-48,tfm-prereq-49,tfm-prereq-5,tfm-prereq-50,tfm-prereq-51,tfm-prereq-52,tfm-prereq-53,tfm-prereq-54,tfm-prereq-55,tfm-prereq-56,tfm-prereq-57,tfm-prereq-58,tfm-prereq-59,tfm-prereq-6,tfm-prereq-60,tfm-prereq-61,tfm-prereq-62,tfm-prereq-63,tfm-prereq-64,tfm-prereq-65,tfm-prereq-66,tfm-prereq-67,tfm-prereq-68,tfm-prereq-69,tfm-prereq-7,tfm-prereq-70,tfm-prereq-71,tfm-prereq-72,tfm-prereq-73,tfm-prereq-74,tfm-prereq-75,tfm-prereq-76,tfm-prereq-77,tfm-prereq-78,tfm-prereq-79,tfm-prereq-8,tfm-prereq-80,tfm-prereq-81,tfm-prereq-82,tfm-prereq-83,tfm-prereq-84,tfm-prereq-85,tfm-prereq-86,tfm-prereq-87,tfm-prereq-88,tfm-prereq-89,tfm-prereq-9,tfm-prereq-90,tfm-prereq-91,tfm-prereq-92,tfm-prereq-93,tfm-prereq-94,tfm-prereq-95,tfm-prereq-96,tfm-prereq-97,tfm-prereq-98,tfm-prereq-99"

for ORG in $(echo $SRC_ORGS | tr ',' ' ')
do
  echo "Migrating organization: $ORG"

  export SRC_TFE_ORG=$ORG

    ./tfm copy workspaces --autoapprove --create-dst-project true 

  echo "Finalized"
done


