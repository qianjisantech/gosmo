import {PaginationRequest} from "@/common/types";

export interface QueryForm extends PaginationRequest{
    planNo:string
    url:string
}