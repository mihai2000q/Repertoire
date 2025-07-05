import { ReactNode } from 'react'
import { Drawer, ScrollArea } from '@mantine/core'
import useTitleBarHeight from '../../../hooks/useTitleBarHeight.ts'

interface RightSideEntityDrawerProps {
  opened: boolean
  onClose: () => void
  isLoading: boolean
  loader: ReactNode
  withScrollArea?: boolean
  children?: ReactNode
}

function RightSideEntityDrawer({
  opened,
  onClose,
  isLoading,
  loader,
  withScrollArea,
  children
}: RightSideEntityDrawerProps) {
  const titleBarHeight = useTitleBarHeight()

  return (
    <Drawer.Root
      opened={opened}
      onClose={onClose}
      position="right"
      size={'max(28%, 440px)'}
      shadow="xl"
      radius={'8 0 0 8'}
      styles={{
        overlay: {
          height: `calc(100% - ${titleBarHeight})`,
          marginTop: titleBarHeight
        },
        inner: {
          marginTop: titleBarHeight
        },
        body: {
          padding: 0,
          margin: 0
        }
      }}
    >
      <Drawer.Overlay backgroundOpacity={0.1} blur={1} />
      <Drawer.Content>
        <Drawer.Body>
          {withScrollArea === false ? (
            <div>{isLoading ? loader : children && children}</div>
          ) : (
            <ScrollArea.Autosize
              mah={`calc(100vh - ${titleBarHeight})`}
              scrollbars={'y'}
              scrollbarSize={10}
              styles={{
                viewport: {
                  '> div': {
                    minWidth: '100%',
                    width: 0
                  }
                }
              }}
            >
              {isLoading ? loader : children && children}
            </ScrollArea.Autosize>
          )}
        </Drawer.Body>
      </Drawer.Content>
    </Drawer.Root>
  )
}

export default RightSideEntityDrawer
