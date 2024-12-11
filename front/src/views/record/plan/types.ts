export interface CreatePlanForm {
    executionMode: string
    project: string
    machine: string
    duration: number
}
export interface Option {
    id: string
    label: string
    value: string
}

export interface MatchTableColumn {
    id: string
    matchMode: string
    matchValue: string
}