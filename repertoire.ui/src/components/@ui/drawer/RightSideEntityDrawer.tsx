import { ReactNode } from 'react'
import { Drawer } from '@mantine/core'
import useTitleBarHeight from '../../../hooks/useTitleBarHeight.ts'

interface RightSideEntityDrawerProps {
  opened: boolean
  onClose: () => void
  isLoading: boolean
  loader: ReactNode
  children?: ReactNode
}

function RightSideEntityDrawer({
  opened,
  onClose,
  isLoading,
  loader,
  children
}: RightSideEntityDrawerProps) {
  const titleBarHeight = useTitleBarHeight()

  return (
    <Drawer
      withCloseButton={false}
      opened={opened}
      onClose={onClose}
      position="right"
      overlayProps={{ backgroundOpacity: 0.1, blur: 1 }}
      size={'max(28%, 440px)'}
      shadow="xl"
      radius={'8 0 0 8'}
      styles={{
        overlay: {
          height: `calc(100% - ${titleBarHeight})`,
          marginTop: titleBarHeight
        },
        inner: {
          height: `calc(100% - ${titleBarHeight})`,
          marginTop: titleBarHeight
        },
        body: {
          padding: 0,
          margin: 0
        }
      }}
    >
      {isLoading ? loader : children && children}
    </Drawer>
  )
}

export default RightSideEntityDrawer
