import Order from "../types/Order.ts";
import OrderType from "./enums/OrderType.ts";

export default function createOrder(order: Order): string {
  return (
    `${order.property} ` +
    `${order.type ?? OrderType.Ascending}` +
    `${
      order.nullable === true
        ? ' ' + (order.type === OrderType.Descending ? 'nulls last' : 'nulls first')
        : ''
    }`
  )
}
