import { mantineRender } from '../../test-utils'
import NewSongCard from './NewSongCard'
import { userEvent } from '@testing-library/user-event'

describe('New Song Card', () => {
  it('should render and when clicked should open the modal', async ({ expect }) => {
    // Arrange
    const openModal = vi.fn()

    const user = userEvent.setup()

    // Act
    const { container } = mantineRender(<NewSongCard openModal={openModal} />)

    // Assert
    const icon = container.querySelector('svg.tabler-icon-music-plus')
    expect(icon).toBeVisible()

    await user.click(icon)

    expect(openModal).toHaveBeenCalledOnce()
  })
})
