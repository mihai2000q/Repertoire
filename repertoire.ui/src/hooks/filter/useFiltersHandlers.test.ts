import { act, renderHook } from '@testing-library/react'
import useFiltersHandlers from './useFiltersHandlers.ts'
import Filter from '../../types/Filter.ts'
import FilterOperator from '../../types/enums/FilterOperator.ts'

describe('useFiltersHandlers', () => {
  const initialFilters = new Map<string, Filter>([
    ['name=', { property: 'name', operator: FilterOperator.Equal, value: '', isSet: false }],
    ['age>', { property: 'age', operator: FilterOperator.GreaterThan, value: 0, isSet: false }],
    [
      'dateFrom>',
      { property: 'dateFrom', operator: FilterOperator.GreaterThan, value: null, isSet: false }
    ],
    [
      'dateTo<',
      { property: 'dateTo', operator: FilterOperator.LessThan, value: null, isSet: false }
    ]
  ])

  const render = (initialState: Map<string, Filter> = new Map<string, Filter>()) => {
    let filters = initialState
    const setFilters = vi.fn((newFilters) => {
      if (typeof newFilters === 'function') {
        filters = newFilters(filters)
      } else {
        filters = newFilters
      }
    })

    const { result } = renderHook(() => useFiltersHandlers(filters, setFilters, initialFilters))
    return { result, setFilters, filters }
  }

  describe('handle IsSet Change', () => {
    it('should update isSet value for the specified filter', () => {
      const { result, setFilters } = render(initialFilters)

      act(() => result.current.handleIsSetChange('name=', true))

      expect(setFilters).toHaveBeenCalledTimes(1)
      const updatedFilters = setFilters.mock.calls[0][0]
      expect(updatedFilters.get('name=')).toStrictEqual({
        ...initialFilters.get('name='),
        isSet: true
      })
    })

    it('should not modify other filters', () => {
      const { result, setFilters } = render(initialFilters)

      act(() => result.current.handleIsSetChange('name=', true))

      const updatedFilters = setFilters.mock.calls[0][0]
      expect(updatedFilters.get('age>')).toStrictEqual(initialFilters.get('age>'))
    })
  })

  describe('handle Value Change', () => {
    it('should update value and set isSet to true when value differs from initial', () => {
      const { result, setFilters } = render(initialFilters)

      act(() => result.current.handleValueChange('name=', 'John'))

      const updatedFilters = setFilters.mock.calls[0][0]
      expect(updatedFilters.get('name=')).toStrictEqual({
        ...initialFilters.get('name='),
        value: 'John',
        isSet: true
      })
    })

    it('should set isSet to false when value matches initial', () => {
      const initialState = new Map<string, Filter>([
        ['name=', { ...initialFilters.get('name='), value: 'John', isSet: true }]
      ])
      const { result, setFilters } = render(initialState)

      act(() => result.current.handleValueChange('name=', ''))

      const updatedFilters = setFilters.mock.calls[0][0]
      expect(updatedFilters.get('name=').isSet).toBeFalsy()
    })

    it('should handle null/undefined/empty string values correctly', () => {
      const { result, setFilters } = render(initialFilters)

      act(() => result.current.handleValueChange('name=', null))

      const updatedFilters = setFilters.mock.calls[0][0]
      expect(updatedFilters.get('name=').isSet).toBeFalsy()
    })
  })

  describe('handle Double Value Change', () => {
    it('should update both values and isSet flags correctly', () => {
      const { result, setFilters } = render(initialFilters)

      act(() =>
        result.current.handleDoubleValueChange('dateFrom>', '2023-01-01', 'dateTo<', '2023-01-31')
      )

      const updatedFilters = setFilters.mock.calls[0][0]
      expect(updatedFilters.get('dateFrom>')).toStrictEqual({
        ...initialFilters.get('dateFrom>'),
        value: '2023-01-01',
        isSet: true
      })
      expect(updatedFilters.get('dateTo<')).toStrictEqual({
        ...initialFilters.get('dateTo<'),
        value: '2023-01-31',
        isSet: true
      })
    })

    it('should handle when one value matches initial and other does not', () => {
      const { result, setFilters } = render(initialFilters)

      act(() => result.current.handleDoubleValueChange('dateFrom>', null, 'dateTo<', '2023-01-31'))

      const updatedFilters = setFilters.mock.calls[0][0]
      expect(updatedFilters.get('dateFrom>').isSet).toBeFalsy()
      expect(updatedFilters.get('dateTo<').isSet).toBeTruthy()
    })
  })

  describe('get Date Range Values', () => {
    it('should return Date objects for valid date strings', () => {
      const initialState = new Map<string, Filter>([
        ['dateFrom>', { ...initialFilters.get('dateFrom>'), value: '2023-01-01' }],
        ['dateTo<', { ...initialFilters.get('dateTo<'), value: '2023-01-31' }]
      ])
      const { result } = render(initialState)

      const [date1, date2] = result.current.getDateRangeValues('dateFrom>', 'dateTo<')
      expect(date1).toBeInstanceOf(Date)
      expect(date2).toBeInstanceOf(Date)
      expect(date1.toISOString()).toMatch(/2023-01-01/)
      expect(date2.toISOString()).toMatch(/2023-01-31/)
    })

    it('should return null for null values', () => {
      const { result } = render(initialFilters)

      const [date1, date2] = result.current.getDateRangeValues('dateFrom>', 'dateTo<')
      expect(date1).toBeNull()
      expect(date2).toBeNull()
    })

    it('should adjust dates when they are the same', () => {
      const sameDate = '2023-01-01'
      const initialState = new Map<string, Filter>([
        ['dateFrom>', { ...initialFilters.get('dateFrom>'), value: sameDate }],
        ['dateTo<', { ...initialFilters.get('dateTo<'), value: sameDate }]
      ])
      const { result } = render(initialState)

      const [date1, date2] = result.current.getDateRangeValues('dateFrom>', 'dateTo<')
      expect(date1.getDate()).toBe(1)
      expect(date2.getDate()).toBe(2)
    })
  })
})
