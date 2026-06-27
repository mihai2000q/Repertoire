import Order from '../../types/Order.ts'
import { IconArrowNarrowDown, IconArrowNarrowUp } from '@tabler/icons-react'
import {
  ActionIcon,
  alpha,
  Center,
  Group,
  Menu,
  ScrollArea,
  Stack,
  Text,
  Tooltip
} from '@mantine/core'
import { MouseEvent, ReactElement, ReactNode, useState } from 'react'
import OrderType from '../../types/enums/OrderType.ts'
import { DragDropContext, Draggable, Droppable } from '@hello-pangea/dnd'
import fixDraggableInMenu from '../../utils/fixDraggableInMenu.ts'

interface AdvancedOrderMenuProps {
  children: ReactNode
  orders: Order[]
  setOrders: (orders: Order[]) => void
  propertyIcons?: Map<string, ReactElement>
}

function AdvancedOrderMenu({ children, orders, setOrders, propertyIcons }: AdvancedOrderMenuProps) {
  const [menuOpened, setMenuOpened] = useState(false)

  const isOnlyOneOrderActive = orders.filter((o) => o.checked === true).length === 1

  function onClick(order: Order, index: number) {
    if (order.checked && isOnlyOneOrderActive) return
    const newOrders = [...orders]
    newOrders[index] = { ...order, checked: !(order.checked ?? false) }
    setOrders(newOrders)
  }

  function onTypeUpdate(e: MouseEvent, order: Order, index: number) {
    e.stopPropagation()
    const newOrders = [...orders]
    newOrders[index] = {
      ...order,
      type: order.type !== OrderType.Descending ? OrderType.Descending : OrderType.Ascending
    }
    setOrders(newOrders)
  }

  function onDragEnd({ source, destination }) {
    const from = source.index
    const to = destination?.index || 0

    if (from === to) return

    const newOrders = [...orders]
    const item = newOrders[from]
    newOrders.splice(from, 1)
    newOrders.splice(to, 0, item)
    setOrders(newOrders)
  }

  return (
    <Menu opened={menuOpened} onChange={setMenuOpened} position={'bottom'}>
      <Menu.Target>{children}</Menu.Target>

      <Menu.Dropdown>
        <ScrollArea.Autosize mah={'40vh'} scrollbars={'y'} scrollbarSize={7}>
          <DragDropContext onDragEnd={onDragEnd}>
            <Droppable droppableId="dnd-list" direction="vertical">
              {(provided) => (
                <Stack ref={provided.innerRef} {...provided.droppableProps} gap={2}>
                  {orders.map((order, index) => (
                    <Draggable key={order.label} index={index} draggableId={order.label}>
                      {(provided, snapshot) => (
                        <Menu.Item
                          component={'div'}
                          ref={provided.innerRef}
                          aria-label={order.label}
                          data-active={order.checked ?? false}
                          px={12}
                          py={8}
                          {...provided.dragHandleProps}
                          {...provided.draggableProps}
                          style={(theme) => ({
                            ...fixDraggableInMenu(provided.draggableProps, snapshot, 79).style,
                            ...(snapshot.isDragging && {
                              cursor: 'grabbing',
                              borderColor: theme.colors.primary[3]
                            })
                          })}
                          sx={(theme) => ({
                            transition: '0.2s',
                            cursor:
                              order.checked === true && isOnlyOneOrderActive
                                ? 'default'
                                : 'pointer',
                            border: '1px solid transparent',
                            color: theme.colors.gray[6],
                            backgroundColor: alpha(theme.white, 0.7),
                            '&:hover': {
                              color: theme.colors.gray[6],
                              backgroundColor: theme.colors.gray[1]
                            },

                            ...(order.checked && {
                              color: theme.colors.primary[4],
                              backgroundColor: alpha(theme.colors.primary[1], 0.4),
                              '&:hover': {
                                color: isOnlyOneOrderActive
                                  ? theme.colors.primary[4]
                                  : theme.colors.primary[5],
                                backgroundColor: isOnlyOneOrderActive
                                  ? alpha(theme.colors.primary[1], 0.4)
                                  : alpha(theme.colors.primary[1], 0.8)
                              }
                            })
                          })}
                          onClick={() => onClick(order, index)}
                        >
                          <Group gap={0}>
                            <Center w={15} h={15}>
                              {propertyIcons?.get(order.property)}
                            </Center>
                            <Text
                              fz={'sm'}
                              c={order.checked === true ? 'primary.4' : 'dark'}
                              inline
                              pl={8}
                              pr={4}
                              flex={1}
                            >
                              {order.label}
                            </Text>
                            {order.checked === true && (
                              <Tooltip label={'Change Order'} openDelay={200}>
                                <ActionIcon
                                  variant={'subtle'}
                                  size={'xs'}
                                  aria-label={
                                    order.type !== OrderType.Descending
                                      ? 'change-order-ascending'
                                      : 'change-order-descending'
                                  }
                                  onClick={(e) => onTypeUpdate(e, order, index)}
                                >
                                  {order.type !== OrderType.Descending ? (
                                    <IconArrowNarrowUp size={16} />
                                  ) : (
                                    <IconArrowNarrowDown size={16} />
                                  )}
                                </ActionIcon>
                              </Tooltip>
                            )}
                          </Group>
                        </Menu.Item>
                      )}
                    </Draggable>
                  ))}
                  {provided.placeholder}
                </Stack>
              )}
            </Droppable>
          </DragDropContext>
        </ScrollArea.Autosize>
      </Menu.Dropdown>
    </Menu>
  )
}

export default AdvancedOrderMenu
