

export  interface PaginationRequest {
    "pageSize": number
    "currentPage":number
}

export  interface PaginationResponse {
    "pageSize": number
    "currentPage":number
}

export interface SearchDataType {
    request_path?: string;
    creator?: number;
    creator_name?: string;
    date_range?: [string, string];
}

export interface TableDataType {
    id?: number;
    index?: number;
    request_path?: string;
    request_method?: string;
    request_ip?: string;
    request_browser?: string;
    request_os?: string;
    response_code?: number;
    request_payload: string;
    response_json?: string;
    creator?: CreatorTableDataType;
    created_at?: string;
    updated_at?: string;
}

export interface SearchCreatorDataType {
    name?: string;
    available?: string;
}

export interface CreatorTableDataType {
    id?: number;
    name?: string;
    available?: string;
    description?: string;
}