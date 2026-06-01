import type { MenuProps, MenuTargetProps } from '@mantine/core'
import { createEventHandler, createSafeContext, isElement, Menu } from '@mantine/core'
import React, { cloneElement, forwardRef, useCallback, useRef } from 'react'
import { useUncontrolled } from '@mantine/hooks'

// Original Credits to: https://gist.github.com/minosss/f26fae6170d62df26103a0c589bf6da6
// Although, the above version only works on React 18

type TriggerEvent = 'click' | 'context'

interface ContextMenuContext {
  opened: boolean
  trigger?: TriggerEvent
  closeDropdown(): void
  toggleDropdown(e: React.MouseEvent, targetElement: HTMLElement): void
}

const [ContextMenuProvider, useContextMenuContext] = createSafeContext<ContextMenuContext>(
  'ContextMenuContext is undefined'
)

interface RefWrapperProps {
  children: React.ReactElement
  refProp: string
}

const RefWrapper = forwardRef<HTMLElement, RefWrapperProps>((props, ref) => {
  const { children, refProp } = props
  const ctx = useContextMenuContext()
  const elementRef = useRef<HTMLElement | null>(null)

  const childProps = children.props as {
    onContextMenu?: React.MouseEventHandler
    onClick?: React.MouseEventHandler
    ref?: React.Ref<HTMLElement>
  }

  const onContextMenu = createEventHandler(childProps.onContextMenu, (e: React.MouseEvent) => {
    if (ctx.trigger === 'context') {
      e.preventDefault()
      if (elementRef.current) {
        ctx.toggleDropdown(e, elementRef.current)
      }
    }
  })

  const onClick = createEventHandler(childProps.onClick, () => ctx.closeDropdown())

  const handleRef = useCallback(
    (node: HTMLElement | null) => {
      if (!node) return
      elementRef.current = node

      if (typeof ref === 'function') ref(node)
      else if (ref) ref.current = node

      const childRef = childProps.ref
      if (childRef) {
        if (typeof childRef === 'function') childRef(node)
        else if (childRef && 'current' in childRef) childRef.current = node
      }
    },
    [ref, childProps.ref]
  )

  if (!isElement(children)) {
    throw new Error('ContextMenu.Target children must be an element or component that accepts ref')
  }

  const additionalProps = {
    onContextMenu,
    onClick,
    [refProp]: handleRef
  }

  // The child element is guaranteed to accept these props (it's a DOM element or component that forwards them)
  return cloneElement(children as React.ReactElement, additionalProps)
})

RefWrapper.displayName = 'RefWrapper'

const ContextMenuTarget = forwardRef<HTMLElement, MenuTargetProps>((props, ref) => {
  const { children, refProp = 'ref', ...others } = props

  if (!isElement(children)) {
    throw new Error('ContextMenu.Target expects a single child element')
  }

  const targetProps = {
    ...others,
    [refProp]: ref
  } as MenuTargetProps

  return (
    <Menu.Target {...targetProps}>
      <RefWrapper refProp={refProp}>{children}</RefWrapper>
    </Menu.Target>
  )
})

ContextMenuTarget.displayName = 'ContextMenuTarget'

export interface ContextMenuProps extends Omit<MenuProps, 'trigger'> {
  trigger?: TriggerEvent
}

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

  const [_opened, _setOpened] = useUncontrolled({
    value: opened,
    defaultValue: defaultOpened,
    finalValue: false,
    onChange
  })

  const targetElementRef = useRef<HTMLElement | null>(null)
  const originalGetBoundingClientRectRef = useRef<(() => DOMRect) | null>(null)

  const restoreOriginalRect = () => {
    if (targetElementRef.current && originalGetBoundingClientRectRef.current) {
      targetElementRef.current.getBoundingClientRect = originalGetBoundingClientRectRef.current
      originalGetBoundingClientRectRef.current = null
    }
  }

  const _close = () => {
    if (disabled) return
    _setOpened(false)
    if (_opened) onClose?.()
    restoreOriginalRect()
  }

  const _open = () => {
    if (disabled) return
    _setOpened(true)
    if (!_opened) onOpen?.()
  }

  const _toggleDropdown = (e: React.MouseEvent, targetElement: HTMLElement) => {
    if (disabled) return

    targetElementRef.current = targetElement
    originalGetBoundingClientRectRef.current =
      targetElement.getBoundingClientRect.bind(targetElement)

    targetElement.getBoundingClientRect = () => ({
      x: e.clientX,
      y: e.clientY,
      width: 0,
      height: 0,
      top: e.clientY,
      left: e.clientX,
      right: e.clientX,
      bottom: e.clientY,
      toJSON: () => {}
    })

    _setOpened(!_opened)
    if (!_opened) onOpen?.()
    if (_opened) onClose?.()
  }

  const ctx: ContextMenuContext = {
    opened: _opened,
    closeDropdown: _close,
    toggleDropdown: _toggleDropdown,
    trigger: trigger
  }

  return (
    <ContextMenuProvider value={ctx}>
      <Menu
        {...others}
        trigger={trigger === 'context' ? undefined : trigger}
        opened={_opened}
        offset={0}
        onChange={_setOpened}
        onClose={_close}
        onOpen={_open}
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
