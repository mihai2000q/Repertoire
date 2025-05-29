import { emptySong, mantineRender } from '../../../../../test-utils.tsx'
import { Menu } from '@mantine/core'
import OpenLinksMenuItem from './OpenLinksMenuItem.tsx'
import Song from '../../../../../types/models/Song.ts'
import { fireEvent, screen, within } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'

describe('Open Links Menu Item', () => {
  const render = (song: Song, openYoutube: () => void = vi.fn()) =>
    mantineRender(
      <Menu opened={true}>
        <Menu.Dropdown>
          <OpenLinksMenuItem song={song} openYoutube={openYoutube} />
        </Menu.Dropdown>
      </Menu>
    )

  it('should render and be disabled when all links are empty', () => {
    render(emptySong)

    expect(screen.getByRole('menuitem', { name: /open links/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /open links/i })).toBeDisabled()
  })

  it('should render and have working songsterr link', async () => {
    const user = userEvent.setup()
    const song = { ...emptySong, songsterrLink: 'some-link' }

    render(song)

    expect(screen.getByRole('menuitem', { name: /open links/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /open links/i })).not.toBeDisabled()

    await user.hover(screen.getByRole('menuitem', { name: /open links/i }))
    expect(screen.getByRole('menuitem', { name: /songsterr/i })).not.toBeDisabled()
    expect(screen.getByRole('menuitem', { name: /youtube/i })).toBeDisabled()

    expect(
      within(screen.getByRole('link', { name: /songsterr/i })).getByRole('menuitem', {
        name: /songsterr/i
      })
    ).toBeInTheDocument()
    expect(screen.getByRole('link', { name: /songsterr/i })).toBeExternalLink(song.songsterrLink)
  })

  it('should render and have working youtube link', async () => {
    const user = userEvent.setup()
    const song = { ...emptySong, youtubeLink: 'some-link' }
    const openYoutube = vi.fn()

    render(song, openYoutube)

    expect(screen.getByRole('menuitem', { name: /open links/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /open links/i })).not.toBeDisabled()

    await user.hover(screen.getByRole('menuitem', { name: /open links/i }))
    expect(screen.getByRole('menuitem', { name: /songsterr/i })).toBeDisabled()
    expect(screen.getByRole('menuitem', { name: /youtube/i })).not.toBeDisabled()

    fireEvent.click(screen.getByRole('menuitem', { name: /youtube/i }))
    expect(openYoutube).toHaveBeenCalledOnce()
  })
})
