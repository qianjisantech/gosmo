import {reactive} from "vue";
import {PaginationRequest} from "@/common/types";


export  const paginationRequest  =reactive<PaginationRequest>( {
    "pageSize": 10,
    "currentPage": 1,
})

export  const paginationResponse = {
    "pageSize": 10,
    "currentPage": 1,
}