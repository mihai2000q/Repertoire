import Order from '../types/Order.ts'
import createOrder from './createOrder.ts'
import OrderType from './enums/OrderType.ts'

describe('create Order', () => {
  it.each([
    [
      {
        property: 'prop'
      },
      'prop asc'
    ],
    [
      {
        property: 'prop',
        type: OrderType.Ascending
      },
      'prop asc'
    ],
    [
      {
        property: 'prop',
        type: OrderType.Descending
      },
      'prop desc'
    ],
    [
      {
        property: 'prop',
        type: OrderType.Descending,
        nullable: false
      },
      'prop desc'
    ],
    [
      {
        property: 'prop',
        nullable: true
      },
      'prop asc nulls first'
    ],
    [
      {
        property: 'prop',
        type: OrderType.Descending,
        nullable: true
      },
      'prop desc nulls last'
    ]
  ])('should create order', (order: Order, expected) => {
    expect(createOrder(order)).toBe(expected)
  })
})
