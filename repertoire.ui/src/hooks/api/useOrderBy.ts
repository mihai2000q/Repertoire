import Order from '../../types/Order.ts'
import { useState } from 'react'
import createOrder from '../../utils/createOrder.ts'
import useDidUpdateShallow from '../useDidUpdateShallow.ts'

export default function useOrderBy(orders: Order[] | null = null): string[] {
  const [orderBy, setOrderBy] = useState<string[]>(createOrders(orders))
  useDidUpdateShallow(() => setOrderBy(createOrders(orders)), [orders])
  return orderBy
}

function createOrders(orders: Order[] | null): string[] | null {
  const newOrders: string[] = []
  if (orders === null) {
    return newOrders
  } else {
    for (const o of orders.filter((order) => order.checked !== false)) {
      newOrders.push(createOrder(o))
      newOrders.push(...(o.thenBy?.map(createOrder) ?? []))
    }
  }
  return newOrders
}
