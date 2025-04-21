import Filter, { FilterValue } from '../../types/Filter.ts'
import { Dispatch, SetStateAction } from 'react'

export default function useFiltersHandlers(
  filters: Map<string, Filter>,
  setFilters: Dispatch<SetStateAction<Map<string, Filter>>>,
  initialFilters: Map<string, Filter>
) {
  function handleIsSetChange(key: string, isSet: boolean) {
    const newFilters = new Map<string, Filter>([...filters])
    newFilters.set(key, { ...newFilters.get(key), isSet: isSet })
    setFilters(newFilters)
  }

  function handleValueChange(key: string, value: FilterValue) {
    const newFilters = new Map<string, Filter>([...filters])
    const initialFilter = initialFilters.get(key)
    newFilters.set(key, {
      ...newFilters.get(key),
      isSet: initialFilter.value !== value && value !== '' && value !== null && value !== undefined,
      value: value
    })
    setFilters(newFilters)
  }

  function handleDoubleValueChange(
    key1: string,
    value1: FilterValue,
    key2: string,
    value2: FilterValue
  ) {
    const newFilters = new Map<string, Filter>([...filters])

    const initialFilter1 = initialFilters.get(key1)
    newFilters.set(key1, {
      ...newFilters.get(key1),
      isSet:
        initialFilter1.value !== value1 && value1 !== '' && value1 !== null && value1 !== undefined,
      value: value1
    })

    const initialFilter2 = initialFilters.get(key2)
    newFilters.set(key2, {
      ...newFilters.get(key2),
      isSet:
        initialFilter2.value !== value2 && value2 !== '' && value2 !== null && value2 !== undefined,
      value: value2
    })

    setFilters(newFilters)
  }

  function getDateRangeValues(key1: string, key2: string): [Date | null, Date | null] {
    const value1 = filters.get(key1).value
    const value2 = filters.get(key2).value

    const date1 = value1 ? new Date(value1 as string) : null
    const date2 = value2 ? new Date(value2 as string) : null

    if (date1 !== null && date2 !== null && date1.toDateString() === date2.toDateString()) {
      date2.setDate(date2.getDate() + 1)
    }

    return [date1, date2]
  }

  return { handleIsSetChange, handleValueChange, handleDoubleValueChange, getDateRangeValues }
}
