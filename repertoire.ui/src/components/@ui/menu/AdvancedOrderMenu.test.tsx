import Order from '../../../types/Order.ts'
import AlbumProperty from '../../../types/enums/AlbumProperty.ts'
import { mantineRender } from '../../../test-utils.tsx'
import AdvancedOrderMenu from './AdvancedOrderMenu.tsx'
import { ReactNode } from 'react'
import { userEvent } from '@testing-library/user-event'
import { screen, within } from '@testing-library/react'
import OrderType from '../../../types/enums/OrderType.ts'

describe('Advanced Order Menu', () => {
  const initialOrders: Order[] = [
    { label: 'Title', property: AlbumProperty.Title, checked: false, type: OrderType.Descending },
    { label: 'Release Date', property: AlbumProperty.ReleaseDate, nullable: true, checked: false },
    { label: 'Artist', property: AlbumProperty.Artist, checked: true },
    { label: 'Confidence', property: AlbumProperty.Confidence, type: OrderType.Ascending }
  ]

  const childrenTestId = 'test-id'
  const render = ({
    orders,
    setOrders,
    children
  }: { orders?: Order[]; setOrders?: (orders: Order[]) => void; children?: ReactNode } = {}) =>
    mantineRender(
      <AdvancedOrderMenu orders={orders ?? initialOrders} setOrders={setOrders}>
        {children ?? <div data-testid={childrenTestId}>something</div>}
      </AdvancedOrderMenu>
    )

  it('should render', async () => {
    const user = userEvent.setup()

    const testId = 'children-test-id'
    const children = <div data-testid={testId}>something</div>

    render({ children })

    await user.click(screen.getByTestId(testId))
    expect(await screen.findByRole('menu')).toBeInTheDocument()
    initialOrders.forEach((order) => {
      expect(screen.getByRole('button', { name: order.label })).toBeInTheDocument()
      expect(screen.getByRole('button', { name: order.label })).toHaveAttribute(
        'data-active',
        order.checked === true ? 'true' : 'false'
      )
      if (order.checked === true) {
        if (order.type !== OrderType.Descending)
          expect(
            within(screen.getByRole('button', { name: order.label })).getByRole('button', {
              name: 'change-order-ascending'
            })
          ).toBeInTheDocument()
        else
          expect(
            within(screen.getByRole('button', { name: order.label })).getByRole('button', {
              name: 'change-order-descending'
            })
          ).toBeInTheDocument()
      }
    })
  })

  it('should add another order on click', async () => {
    const user = userEvent.setup()

    const setOrders = vi.fn()

    const orderIndex = 0
    const newOrder = initialOrders[orderIndex]

    render({ setOrders })

    await user.click(screen.getByTestId(childrenTestId))
    await user.click(screen.getByRole('button', { name: newOrder.label }))

    const newOrders = [...initialOrders]
    newOrders[orderIndex] = { ...newOrder, checked: true }
    expect(setOrders).toHaveBeenCalledExactlyOnceWith(newOrders)
  })

  it('should not remove orders if only one is active', async () => {
    const user = userEvent.setup()

    const setOrders = vi.fn()

    const order = initialOrders[2]

    render({ setOrders })

    await user.click(screen.getByTestId(childrenTestId))
    await user.click(screen.getByRole('button', { name: order.label }))

    expect(setOrders).not.toHaveBeenCalled()
  })

  it('should remove order on click when there is more than 1 active', async () => {
    const user = userEvent.setup()

    const setOrders = vi.fn()

    const orderIndex = 0
    const orderToRemove = initialOrders[orderIndex]
    const orders = [...initialOrders]
    orders[orderIndex] = { ...orderToRemove, checked: true }

    render({ orders, setOrders })

    await user.click(screen.getByTestId(childrenTestId))
    await user.click(screen.getByRole('button', { name: orderToRemove.label }))

    const newOrders = [...initialOrders]
    newOrders[orderIndex] = { ...orderToRemove, checked: false }
    expect(setOrders).toHaveBeenCalledOnce()
    expect(setOrders).toHaveBeenCalledWith(newOrders)
  })

  it('should change the type of order', async () => {
    const user = userEvent.setup()

    const setOrders = vi.fn()

    const orderIndex = 2
    const order = initialOrders[orderIndex]

    render({ setOrders })

    await user.click(screen.getByTestId(childrenTestId))
    await user.click(
      within(screen.getByRole('button', { name: order.label })).getByRole('button', {
        name: 'change-order-ascending'
      })
    )

    const newOrders = [...initialOrders]
    newOrders[orderIndex] = { ...order, type: OrderType.Descending }
    expect(setOrders).toHaveBeenCalledOnce()
    expect(setOrders).toHaveBeenCalledWith(newOrders)
  })

  it('should be able to reorder orders', async () => {})
})
