import { mantineRender } from '../../../test-utils.tsx'
import HomeTopEntityOrderButton from './HomeTopEntityOrderButton.tsx'
import HomeTopEntity from '../../../types/enums/HomeTopEntity.ts'
import Order from '../../../types/Order.ts'
import { userEvent } from '@testing-library/user-event'
import { screen } from '@testing-library/react'
import homeTopOrderEntities, {
  defaultHomeTopOrderEntities
} from '../../../data/home/homeTopOrderEntities.ts'

describe('Home Top Entity Order Button', () => {
  it('should render', async () => {
    const user = userEvent.setup()

    const topEntity = HomeTopEntity.Albums
    const orderEntities = defaultHomeTopOrderEntities

    mantineRender(
      <HomeTopEntityOrderButton
        topEntity={topEntity}
        orderEntities={orderEntities}
        setOrderEntities={vi.fn()}
      />
    )

    expect(screen.getByRole('button', { name: 'order' })).toBeInTheDocument()

    await user.click(screen.getByRole('button', { name: 'order' }))
    for (const order of homeTopOrderEntities.get(topEntity)) {
      expect(await screen.findByRole('menuitem', { name: order.label })).toBeInTheDocument()
      if (order.property === orderEntities.get(topEntity).property)
        expect(await screen.findByRole('menuitem', { name: order.label })).toHaveAttribute(
          'aria-selected',
          'true'
        )
    }
  })

  it("should render with all entities' order options", async () => {
    const user = userEvent.setup()

    const orderEntities = defaultHomeTopOrderEntities

    const { rerender } = mantineRender(
      <HomeTopEntityOrderButton
        topEntity={HomeTopEntity.Artists}
        orderEntities={defaultHomeTopOrderEntities}
        setOrderEntities={vi.fn()}
      />
    )

    await user.click(screen.getByRole('button', { name: 'order' }))
    for (const [key, orders] of homeTopOrderEntities) {
      rerender(
        <HomeTopEntityOrderButton
          topEntity={key}
          orderEntities={defaultHomeTopOrderEntities}
          setOrderEntities={vi.fn()}
        />
      )

      for (const order of orders) {
        expect(await screen.findByRole('menuitem', { name: order.label })).toBeInTheDocument()
        if (order.property === orderEntities.get(key).property)
          expect(await screen.findByRole('menuitem', { name: order.label })).toHaveAttribute(
            'aria-selected',
            'true'
          )
      }
    }
  })

  it('should call set order entities when clicking on new order', async () => {
    const user = userEvent.setup()

    const topEntity = HomeTopEntity.Albums
    const orderEntities = new Map<HomeTopEntity, Order>([
      [topEntity, homeTopOrderEntities.get(topEntity)[0]],
      [HomeTopEntity.Songs, homeTopOrderEntities.get(HomeTopEntity.Songs)[0]]
    ])
    const setOrderEntities = vi.fn()

    const newOrder = homeTopOrderEntities.get(topEntity)[1]

    mantineRender(
      <HomeTopEntityOrderButton
        topEntity={topEntity}
        orderEntities={orderEntities}
        setOrderEntities={setOrderEntities}
      />
    )

    expect(screen.getByRole('button', { name: 'order' })).toBeInTheDocument()

    await user.click(screen.getByRole('button', { name: 'order' }))
    await user.click(await screen.findByRole('menuitem', { name: newOrder.label }))

    expect(setOrderEntities).toHaveBeenCalledExactlyOnceWith(expect.any(Function))

    const newOrderEntities = setOrderEntities.mock.calls[0][0](orderEntities) as Map<
      HomeTopEntity,
      Order
    >
    expect(newOrderEntities).toHaveLength(orderEntities.size)
    for (const [key, value] of newOrderEntities) {
      if (key === topEntity) expect(JSON.stringify(value)).toBe(JSON.stringify(newOrder))
      else expect(JSON.stringify(value)).toBe(JSON.stringify(orderEntities.get(key)))
    }
  })
})
