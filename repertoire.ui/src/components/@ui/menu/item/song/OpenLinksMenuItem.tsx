import { Anchor, Center, Menu } from '@mantine/core'
import Song from '../../../../../types/models/Song.ts'
import { IconBrandYoutubeFilled, IconGuitarPickFilled, IconLocation } from '@tabler/icons-react'

interface OpenLinksMenuItemProps {
  song: Song
  openYoutube: () => void
}

function OpenLinksMenuItem({ song, openYoutube }: OpenLinksMenuItemProps) {
  return (
    <Menu.Sub position={'right'}>
      <Menu.Sub.Target>
        <Menu.Sub.Item
          leftSection={<IconLocation size={14} />}
          disabled={!song.songsterrLink && !song.youtubeLink}
        >
          Open Links
        </Menu.Sub.Item>
      </Menu.Sub.Target>

      <Menu.Sub.Dropdown>
        <Anchor
          underline={'never'}
          href={song.songsterrLink}
          target="_blank"
          rel="noreferrer"
          c={'inherit'}
          onClick={(e) => e.stopPropagation()}
        >
          <Menu.Item
            leftSection={
              <Center c={'blue.7'}>
                <IconGuitarPickFilled size={14} />
              </Center>
            }
            disabled={!song.songsterrLink}
          >
            Open Songsterr
          </Menu.Item>
        </Anchor>

        <Menu.Item
          leftSection={
            <Center c={'red.7'}>
              <IconBrandYoutubeFilled size={14} />
            </Center>
          }
          disabled={!song.youtubeLink}
          onClick={openYoutube}
        >
          Open Youtube
        </Menu.Item>
      </Menu.Sub.Dropdown>
    </Menu.Sub>
  )
}

export default OpenLinksMenuItem
