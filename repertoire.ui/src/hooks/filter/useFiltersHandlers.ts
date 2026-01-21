import Filter, { FilterValue } from '../../types/Filter.ts'
import { Dispatch, SetStateAction } from 'react'

export default function useFiltersHandlers(
  filters: Map<string, Filter>,
  setFilters: Dispatch<SetStateAction<Map<string, Filter>>>,
  initialFilters?: Map<string, Filter>
) {
  function handleIsSetChange(key: string, isSet: boolean) {
    const newFilters = new Map<string, Filter>([...filters])
    newFilters.set(key, { ...newFilters.get(key), isSet: isSet })
    setFilters(newFilters)
  }

  function handleValueChange(key: string, value: FilterValue) {
    const newFilters = new Map<string, Filter>([...filters])
    const initialFilter = initialFilters?.get(key)
    newFilters.set(key, {
      ...newFilters.get(key),
      isSet: isSet(value, initialFilter),
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

    const initialFilter1 = initialFilters?.get(key1)
    newFilters.set(key1, {
      ...newFilters.get(key1),
      isSet: isSet(value1, initialFilter1),
      value: value1
    })

    const initialFilter2 = initialFilters?.get(key2)
    newFilters.set(key2, {
      ...newFilters.get(key2),
      isSet: isSet(value2, initialFilter2),
      value: value2
    })

    setFilters(newFilters)
  }

  function getDateRangeValues(key1: string, key2: string): [string | null, string | null] {
    const value1 = filters.get(key1).value
    const value2 = filters.get(key2).value

    const date1 = value1 ? new Date(value1 as string) : null
    const date2 = value2 ? new Date(value2 as string) : null

    if (date1 !== null && date2 !== null && date1.toDateString() === date2.toDateString()) {
      date2.setDate(date2.getDate() + 1)
    }

    return [date1?.toISOString() ?? null, date2?.toISOString() ?? null]
  }

  function getSliderValues(key1: string, key2: string): [number, number] {
    const value1 = filters.get(key1).value
    const value2 = filters.get(key2).value

    let number1 = (value1 ?? 0) as number
    let number2 = (value2 ?? 0) as number

    if (number1 === number2 && number2 === 100) {
      number1--
    } else if (value1 === value2) {
      number2++
    }

    return [number1, number2]
  }

  return {
    handleIsSetChange,
    handleValueChange,
    handleDoubleValueChange,
    getDateRangeValues,
    getSliderValues
  }
}

function isSet(value: FilterValue, initialFilter: Filter): boolean {
  return (
    (!initialFilter || initialFilter.value !== value) &&
    value !== '' &&
    (!Array.isArray(value) || value.length !== 0) &&
    value !== null &&
    value !== undefined
  )
}
