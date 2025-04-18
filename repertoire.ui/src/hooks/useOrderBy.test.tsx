import { act, renderHook } from '@testing-library/react'
import useOrderBy from './useOrderBy'
import Order from '../types/Order'
import OrderType from '../types/enums/OrderType'

describe('use Order By', () => {
  it('should return empty array for null orders', () => {
    const { result } = renderHook(() => useOrderBy(null))
    expect(result.current).toStrictEqual([])
  })

  it('should handle single order with default type', () => {
    const orders: Order[] = [{ property: 'name', checked: true }]
    const { result } = renderHook(() => useOrderBy(orders))
    expect(result.current).toStrictEqual(['name asc'])
  })

  it('should handle explicit order type', () => {
    const orders: Order[] = [{ property: 'name', type: OrderType.Descending, checked: true }]
    const { result } = renderHook(() => useOrderBy(orders))
    expect(result.current).toStrictEqual(['name desc'])
  })

  it('should handle nullable ascending', () => {
    const orders: Order[] = [{ property: 'name', nullable: true, checked: true }]
    const { result } = renderHook(() => useOrderBy(orders))
    expect(result.current).toStrictEqual(['name asc nulls first'])
  })

  it('should handle nullable descending', () => {
    const orders: Order[] = [
      { property: 'name', type: OrderType.Descending, nullable: true, checked: true }
    ]
    const { result } = renderHook(() => useOrderBy(orders))
    expect(result.current).toStrictEqual(['name desc nulls last'])
  })

  it('should ignore unchecked orders', () => {
    const orders: Order[] = [
      { property: 'name', checked: false },
      { property: 'age', checked: true }
    ]
    const { result } = renderHook(() => useOrderBy(orders))
    expect(result.current).toStrictEqual(['age asc'])
  })

  it('should handle thenBy orders', () => {
    const orders: Order[] = [
      {
        property: 'department',
        checked: true,
        thenBy: [
          { property: 'team' },
          { property: 'location', type: OrderType.Descending },
          { property: 'month', type: OrderType.Descending, nullable: true }
        ]
      }
    ]

    const { result } = renderHook(() => useOrderBy(orders))
    expect(result.current).toStrictEqual([
      'department asc',
      'team asc',
      'location desc',
      'month desc nulls last'
    ])
  })

  it('should update when orders change', () => {
    const initialOrders: Order[] = [{ property: 'name', checked: true }]

    const { result, rerender } = renderHook(({ orders }) => useOrderBy(orders), {
      initialProps: { orders: initialOrders }
    })

    expect(result.current).toStrictEqual(['name asc'])

    const newOrders = [
      { property: 'date', type: OrderType.Descending, nullable: true, checked: true }
    ]
    act(() => rerender({ orders: newOrders }))

    expect(result.current).toStrictEqual(['date desc nulls last'])
  })

  it('should handle complex mixed cases', () => {
    const orders: Order[] = [
      {
        property: 'company',
        nullable: true,
        checked: true,
        thenBy: [{ property: 'department', type: OrderType.Descending }]
      },
      {
        property: 'salary',
        type: OrderType.Descending,
        checked: false // ignored
      },
      { property: 'hireDate', checked: true },
      { property: 'salary', type: OrderType.Ascending }
    ]

    const { result } = renderHook(() => useOrderBy(orders))
    expect(result.current).toStrictEqual([
      'company asc nulls first',
      'department desc',
      'hireDate asc',
      'salary asc'
    ])
  })
})
