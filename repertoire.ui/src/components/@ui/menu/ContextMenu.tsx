import type { MenuProps, MenuTargetProps } from '@mantine/core'
import { createEventHandler, createSafeContext, isElement, Menu } from '@mantine/core'
import React, { cloneElement, forwardRef, useRef } from 'react'
import { mergeRefs, useUncontrolled } from '@mantine/hooks'

// Credits to: https://gist.github.com/minosss/f26fae6170d62df26103a0c589bf6da6

type TriggerEvent = 'click' | 'context'

interface ContextMenuContext {
  opened: boolean

  toggleDropdown(e: React.MouseEvent): void

  trigger?: TriggerEvent
}

const [ContextMenuProvider, useContextMenuContext] = createSafeContext<ContextMenuContext>(
  'ContextMenuContext is undefined'
)

type RefWrapperProps = React.PropsWithChildren<{ refProp: string }>

/** ref wrapper, append custom floating middleware to move dropdown follow mouse click */
const RefWrapper = forwardRef<HTMLElement, RefWrapperProps>((props, ref) => {
  const { children, refProp } = props

  if (!isElement(children)) {
    throw new Error(
      'ContextMenu.Target component children should be an element or a component that accepts ref'
    )
  }
  const ctx = useContextMenuContext()

  const toggleDropdown = (e: React.MouseEvent) => {
    // ref of trigger should be a function
    if (typeof ref === 'function') {
      if (!ctx.opened) {
        ref({
          // eslint-disable-next-line @typescript-eslint/ban-ts-comment
          // @ts-ignore
          getBoundingClientRect() {
            return {
              x: e.clientX,
              y: e.clientY,
              width: 0,
              height: 0,
              top: e.clientY,
              right: e.clientX,
              bottom: e.clientY,
              left: e.clientX
            }
          }
        })
      }
      ctx.toggleDropdown(e)
    }
  }

  const onContextMenu = createEventHandler(children.props.onContextMenu, (e: React.MouseEvent) => {
    if (ctx.trigger === 'context') {
      e.preventDefault()
      toggleDropdown(e)
    }
  })

  return cloneElement(children, {
    onContextMenu,
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore
    [refProp]: mergeRefs(ref, children.ref)
  })
})

RefWrapper.displayName = 'RefWrapper'

const ContextMenuTarget = forwardRef<HTMLElement, MenuTargetProps>((props, ref) => {
  const { children, refProp = 'ref', ...others } = props
  return (
    <Menu.Target {...others} refProp={refProp} ref={ref}>
      <RefWrapper refProp={refProp}>{children}</RefWrapper>
    </Menu.Target>
  )
})

ContextMenuTarget.displayName = 'ContextMenuTarget'

export interface ContextMenuProps extends Omit<MenuProps, 'trigger'> {
  trigger?: TriggerEvent
}

/**
 * ContextMenu, Menu Wrapper make the menu(dropdown) follow the mouse click
 *
 * @example
 * ```tsx
 * <ContextMenu position='top-end'>
 *  <ContextMenu.Target>
 *    <Center h={100} bg='teal'>
 *      Right Click
 *    </Center>
 *  </ContextMenu.Target>
 *  <ContextMenu.Dropdown>
 *    <ContextMenu.Item>Undo</ContextMenu.Item>
 *    <ContextMenu.Item>Redo</ContextMenu.Item>
 *  </ContextMenu.Dropdown>
 * </ContextMenu>
 * ```
 */
export const ContextMenu = (props: ContextMenuProps) => {
  const {
    opened,
    defaultOpened,
    onChange,
    onOpen,
    onClose,
    children,
    disabled,
    trigger = 'context',
    position = 'bottom-start',
    shadow = 'lg',
    transitionProps = { transition: 'pop-top-left', duration: 150 },
    ...others
  } = props

  // controlled menu opened state
  const [_opened, setOpened] = useUncontrolled({
    value: opened,
    defaultValue: defaultOpened,
    finalValue: false,
    onChange
  })

  const close = () => {
    if (disabled) return
    setOpened(false)
    if (_opened) onClose?.()
  }

  const open = () => {
    if (disabled) return
    setOpened(true)
    if (!_opened) onOpen?.()
  }

  const lastEventRef = useRef<React.MouseEvent | null>(null)
  const toggleDropdown = (e: React.MouseEvent) => {
    lastEventRef.current = e
    if (_opened) close()
    else open()
  }

  const ctx: ContextMenuContext = {
    opened: _opened,
    toggleDropdown,
    trigger
  }

  return (
    <ContextMenuProvider value={ctx}>
      <Menu
        {...others}
        trigger={trigger === 'context' ? undefined : trigger}
        opened={_opened}
        offset={0}
        onChange={setOpened}
        onClose={close}
        onOpen={open}
        defaultOpened={defaultOpened}
        disabled={disabled}
        // with default values
        position={position}
        shadow={shadow}
        transitionProps={transitionProps}
      >
        {children}
      </Menu>
    </ContextMenuProvider>
  )
}

ContextMenu.Target = ContextMenuTarget
ContextMenu.Dropdown = Menu.Dropdown
ContextMenu.Label = Menu.Label
ContextMenu.Item = Menu.Item
ContextMenu.Divider = Menu.Divider
