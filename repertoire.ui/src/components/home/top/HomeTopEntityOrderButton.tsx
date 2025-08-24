import { IconFilter2 } from '@tabler/icons-react'
import { ActionIcon, Center, Menu } from '@mantine/core'
import HomeTopEntity from '../../../types/enums/HomeTopEntity.ts'
import Order from '../../../types/Order.ts'
import { Dispatch, SetStateAction, useMemo } from 'react'
import homeTopOrderEntities from '../../../data/home/homeTopOrderEntities.ts'
import { artistPropertyIcons } from '../../../data/icons/artistPropertyIcons.tsx'
import { albumPropertyIcons } from '../../../data/icons/albumPropertyIcons.tsx'
import { songPropertyIcons } from '../../../data/icons/songPropertyIcons.tsx'

interface HomeTopEntityOrderButtonProps {
  topEntity: HomeTopEntity
  orderEntities: Map<HomeTopEntity, Order>
  setOrderEntities: Dispatch<SetStateAction<Map<HomeTopEntity, Order>>>
  disabled?: boolean
}

function HomeTopEntityOrderButton({
  topEntity,
  orderEntities,
  setOrderEntities,
  disabled
}: HomeTopEntityOrderButtonProps) {
  const propertyIcons = useMemo(() => {
    switch (topEntity) {
      case HomeTopEntity.Artists:
        return artistPropertyIcons
      case HomeTopEntity.Albums:
        return albumPropertyIcons
      case HomeTopEntity.Songs:
        return songPropertyIcons
    }
  }, [topEntity])

  function handleChange(order: Order) {
    if (isSelected(order)) return
    setOrderEntities(
      (prev) =>
        new Map<HomeTopEntity, Order>(
          [...prev].map(([key, value]) => [key, key === topEntity ? order : value])
        )
    )
  }

  function isSelected(order: Order) {
    return JSON.stringify(orderEntities.get(topEntity)) === JSON.stringify(order)
  }

  return (
    <Menu position={'bottom'} disabled={disabled}>
      <Menu.Target>
        <ActionIcon aria-label={'order'} size={'md'} variant={'grey'} disabled={disabled}>
          <IconFilter2 size={17} />
        </ActionIcon>
      </Menu.Target>

      <Menu.Dropdown display={'flex'} style={{ flexDirection: 'column', gap: '2px' }}>
        {homeTopOrderEntities.get(topEntity).map((order) => (
          <Menu.Item
            key={order.property}
            aria-selected={isSelected(order)}
            leftSection={
              <Center c={isSelected(order) ? 'inherit' : 'gray.6'} w={15} h={15}>
                {propertyIcons?.get(order.property)}
              </Center>
            }
            sx={(theme) => ({
              '&[aria-selected="true"]': {
                cursor: 'default',
                color: theme.colors.primary[4],
                backgroundColor: theme.colors.primary[0]
              }
            })}
            onClick={() => handleChange(order)}
          >
            {order.label}
          </Menu.Item>
        ))}
      </Menu.Dropdown>
    </Menu>
  )
}

export default HomeTopEntityOrderButton
