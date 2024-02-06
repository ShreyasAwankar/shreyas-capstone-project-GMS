variable "project_id" {
  default = "capstone-412502"
}

variable "region" {
  default = "us-central1"
}

variable "functions" {
  type = map(object({
    zip        = string
    name       = string
    trigger    = string
    runtime    = string
    entrypoint = string
  }))

  default = {
    "createbulkupload" : {
      zip        = "capstone-zips/createBulkUpload.zip"
      name       = "createBulkUpload"
      trigger    = "http-trigger"
      runtime    = "go121"
      entrypoint = "CreateBulkUpload"
    },
    "createitem" : {
      zip        = "capstone-zips/createItem.zip"
      name       = "createItem"
      trigger    = "http-trigger"
      runtime    = "go121"
      entrypoint = "CreateItem"
    },
    "deleteitem" : {
      zip        = "capstone-zips/deleteItem.zip"
      name       = "deleteItem"
      trigger    = "http-trigger"
      runtime    = "go121"
      entrypoint = "DeleteItem"
    },
    "getitembyid" : {
      zip        = "capstone-zips/getItemById.zip"
      name       = "getItemById"
      trigger    = "http-trigger"
      runtime    = "go121"
      entrypoint = "GetItemById"
    },
    "listgrocery" : {
      zip        = "capstone-zips/listGrocery.zip"
      name       = "listGrocery"
      trigger    = "http-trigger"
      runtime    = "go121"
      entrypoint = "ListGrocery"
    },
    "login" : {
      zip        = "capstone-zips/login.zip"
      name       = "login"
      trigger    = "http-trigger"
      runtime    = "go121"
      entrypoint = "Login"
    },
    "updateitem" : {
      zip        = "capstone-zips/updateItem.zip"
      name       = "updateItem"
      trigger    = "http-trigger"
      runtime    = "go121"
      entrypoint = "UpdateItem"
    }
  }
}
