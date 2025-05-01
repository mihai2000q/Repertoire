import FilterOperator from './enums/FilterOperator.ts'

export type FilterValue = string | string[] | number | boolean | null | undefined

export default interface Filter {
  property: string
  operator: FilterOperator
  value?: FilterValue,
  isSet?: boolean
}
