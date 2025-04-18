import OrderType from './enums/OrderType.ts'

export default interface Order {
  label?: string
  property: string
  type?: OrderType
  nullable?: boolean
  thenBy?: ThenByOrder[]
  checked?: boolean
}

interface ThenByOrder {
  property: string
  type?: OrderType
  nullable?: boolean
}
