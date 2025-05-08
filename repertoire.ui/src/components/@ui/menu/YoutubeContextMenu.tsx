import { Children, cloneElement, MouseEvent, ReactElement } from 'react'
import { Anchor, Menu } from '@mantine/core'
import useContextMenu from '../../../hooks/useContextMenu.ts'
import { useDisclosure } from '@mantine/hooks'
import YoutubeModal from '../modal/YoutubeModal.tsx'
import { IconLayoutSidebarLeftExpand, IconWorldUpload } from '@tabler/icons-react'

interface YoutubeContextMenuProps {
  children: ReactElement
  title: string
  link: string
  onContextMenu?: (e: MouseEvent) => void
}

function YoutubeContextMenu({ children, title, link, onContextMenu }: YoutubeContextMenuProps) {
  const [openedMenu, menuDropdownProps, { openMenu, closeMenu }] = useContextMenu()

  const [openedYoutube, { open: openYoutube, close: closeYoutube }] = useDisclosure(false)

  function handleOpenYoutube() {
    openYoutube()
  }

  return (
    <>
      {Children.map(children, (child) =>
        cloneElement(child, {
          onContextMenu: (e) => {
            if (onContextMenu) onContextMenu(e)
            openMenu(e)
          }
        })
      )}
      <Menu shadow={'lg'} opened={openedMenu} onClose={closeMenu}>
        <Menu.Target>
          <div style={{ display: 'none' }}></div>
        </Menu.Target>
        <Menu.Dropdown {...menuDropdownProps}>
          <Menu.Item
            leftSection={<IconLayoutSidebarLeftExpand size={14} />}
            onClick={handleOpenYoutube}
          >
            Open Modal
          </Menu.Item>
          <Anchor
            underline={'never'}
            href={link}
            target="_blank"
            rel="noreferrer"
            c={'inherit'}
            onClick={(e) => e.stopPropagation()}
          >
            <Menu.Item leftSection={<IconWorldUpload size={14} />}>Open Browser</Menu.Item>
          </Anchor>
        </Menu.Dropdown>
      </Menu>

      <YoutubeModal title={title} link={link} opened={openedYoutube} onClose={closeYoutube} />
    </>
  )
}

export default YoutubeContextMenu
