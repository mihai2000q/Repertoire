import { Button, Menu } from '@mantine/core'
import { IconCaretDownFilled, IconCheck } from '@tabler/icons-react'
import Order from '../../../types/Order.ts'
import { Dispatch, SetStateAction } from 'react'

interface CompactOrderButtonProps {
  availableOrders: Order[]
  order: Order
  setOrder: Dispatch<SetStateAction<Order>>
  disabledOrders?: Order[]
}

function CompactOrderButton({ availableOrders, order, setOrder, disabledOrders }: CompactOrderButtonProps) {
  return (
    <Menu shadow={'sm'}>
      <Menu.Target>
        <Button
          variant={'subtle'}
          size={'compact-xs'}
          rightSection={<IconCaretDownFilled size={11} />}
          styles={{ section: { marginLeft: 4 } }}
        >
          {order.label}
        </Button>
      </Menu.Target>

      <Menu.Dropdown>
        {availableOrders.map((o) => (
          <Menu.Item
            key={o.value}
            leftSection={order === o && <IconCheck size={12} />}
            disabled={disabledOrders?.includes(o)}
            onClick={() => setOrder(o)}
          >
            {o.label}
          </Menu.Item>
        ))}
      </Menu.Dropdown>
    </Menu>
  )
}

export default CompactOrderButton
