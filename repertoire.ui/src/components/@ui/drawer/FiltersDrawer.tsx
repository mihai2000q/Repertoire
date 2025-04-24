import { Dispatch, ReactNode, SetStateAction, useEffect } from 'react'
import {
  ActionIcon,
  Button,
  Drawer,
  Group,
  LoadingOverlay,
  ScrollArea,
  Space,
  Text
} from '@mantine/core'
import useTitleBarHeight from '../../../hooks/useTitleBarHeight.ts'
import { IconFilterFilled } from '@tabler/icons-react'
import Filter from '../../../types/Filter.ts'

interface FiltersDrawerProps {
  opened: boolean
  onClose: () => void
  filters: Map<string, Filter>
  setFilters: Dispatch<SetStateAction<Map<string, Filter>>>
  internalFilters: Map<string, Filter>
  setInternalFilters: Dispatch<SetStateAction<Map<string, Filter>>>
  initialFilters: Map<string, Filter>
  children?: ReactNode
  isLoading?: boolean
  additionalReset?: () => void
}

function FiltersDrawer({
  opened,
  onClose,
  filters,
  setFilters,
  internalFilters,
  setInternalFilters,
  initialFilters,
  children,
  isLoading,
  additionalReset
}: FiltersDrawerProps) {
  const titleBarHeight = useTitleBarHeight()
  const offset = '120px'
  const drawerOffset = '8px'
  const titleSize = '60px'

  useEffect(() => setInternalFilters(initialFilters), [])

  const disabledApplyFilters = JSON.stringify([...internalFilters]) === JSON.stringify([...filters])
  const disabledResetFilters =
    JSON.stringify([...internalFilters]) === JSON.stringify([...initialFilters]) &&
    JSON.stringify([...filters]) === JSON.stringify([...initialFilters])

  function handleApplyFilters() {
    setFilters(internalFilters)
  }

  function handleResetFilters() {
    setFilters(initialFilters)
    setInternalFilters(initialFilters)
    if (additionalReset) additionalReset()
  }

  return (
    <Drawer.Root
      trapFocus={false}
      opened={opened}
      onClose={onClose}
      position={'right'}
      size={'max(20%, 250px)'}
      shadow={'xxl'}
      radius={'lg'}
      offset={drawerOffset}
      styles={{
        overlay: {
          height: `calc(100% - ${titleBarHeight})`,
          top: `${titleBarHeight}`
        },
        inner: {
          top: `calc(${titleBarHeight} + ${offset})`
        },
        body: {
          padding: 0,
          margin: 0
        }
      }}
    >
      <Drawer.Overlay backgroundOpacity={0} />
      <Drawer.Content>
        <LoadingOverlay visible={isLoading} loaderProps={{ type: 'bars' }} />

        <Drawer.Header style={{ zIndex: 100 }}>
          <Group w={'100%'} gap={'xxs'}>
            <Text fw={700} fz={'lg'} inline>
              Filters
            </Text>
            <ActionIcon
              aria-label={'apply-filters'}
              variant={'subtle'}
              size={'md'}
              disabled={disabledApplyFilters}
              onClick={handleApplyFilters}
            >
              <IconFilterFilled size={15} />
            </ActionIcon>
            <Space flex={1} />
            <Button
              size={'compact-xs'}
              variant={'subtle'}
              disabled={disabledResetFilters}
              onClick={handleResetFilters}
            >
              Reset
            </Button>
            <Drawer.CloseButton />
          </Group>
        </Drawer.Header>
        <Drawer.Body>
          <ScrollArea.Autosize
            mah={`calc(100vh - ${titleBarHeight} - ${offset} - ${drawerOffset}*2 - ${titleSize})`}
            scrollbars={'y'}
            scrollbarSize={10}
          >
            {children && children}
          </ScrollArea.Autosize>
        </Drawer.Body>
      </Drawer.Content>
    </Drawer.Root>
  )
}

export default FiltersDrawer
