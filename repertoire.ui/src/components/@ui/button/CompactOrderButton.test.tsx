import { mantineRender } from '../../../test-utils.tsx'
import CompactOrderButton from './CompactOrderButton.tsx'
import Order from '../../../types/Order.ts'
import { screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import { expect } from 'vitest'

describe('Compact Order Button', () => {
  const availableOrders: Order[] = [
    { label: 'Order 1', property: 'order_1' },
    { label: 'Order 2', property: 'order_2' },
    { label: 'Order 3', property: 'order_3' }
  ]

  it('should render', async () => {
    const user = userEvent.setup()

    const order = availableOrders[0]
    const setOrder = vitest.fn()
    const disabledOrders = [availableOrders[1]]

    const newOrder = availableOrders[2]

    const { rerender } = mantineRender(
      <CompactOrderButton
        availableOrders={availableOrders}
        order={order}
        setOrder={setOrder}
        disabledOrders={disabledOrders}
      />
    )

    expect(screen.getByRole('button', { name: order.label })).toBeInTheDocument()

    await user.click(screen.getByRole('button', { name: order.label }))
    availableOrders.forEach((o) => {
      const screenOrder = screen.getByRole('menuitem', { name: o.label })
      expect(screenOrder).toBeInTheDocument()
      if (disabledOrders.includes(o)) {
        expect(screenOrder).toBeDisabled()
      } else {
        expect(screenOrder).not.toBeDisabled()
      }
    })

    await user.click(screen.getByRole('menuitem', { name: newOrder.label }))

    expect(setOrder).toHaveBeenCalledOnce()
    expect(setOrder).toHaveBeenCalledWith(newOrder)

    rerender(
      <CompactOrderButton
        availableOrders={availableOrders}
        order={newOrder}
        setOrder={setOrder}
        disabledOrders={disabledOrders}
      />
    )

    expect(screen.getByRole('button', { name: newOrder.label })).toBeInTheDocument()
  })
})
