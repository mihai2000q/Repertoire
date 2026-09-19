import { http, HttpResponse } from 'msw'
import { emptySongSection, reduxRender, withToastify } from '../../../../../../../../test-utils.tsx'
import { setupServer } from 'msw/node'
import { screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import DeleteSongSectionModal from './DeleteSongSectionModal.tsx'
import { SongSection } from '../../../../../../../../types/models/Song.ts'

describe('Delete Song Section Modal', () => {
  const section: SongSection = {
    ...emptySongSection,
    id: '1',
    name: 'Solo 1'
  }

  const server = setupServer()

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  it('should render', () => {
    reduxRender(
      <DeleteSongSectionModal opened={true} onClose={vi.fn()} section={section} songId={'song-1'} />
    )

    expect(screen.getByRole('dialog', { name: /delete section/i })).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: /delete section/i })).toBeInTheDocument()
    expect(screen.getByText(/are you sure/i)).toBeInTheDocument()
    expect(screen.getByText(section.name)).toBeInTheDocument()
    expect(
      screen.getByRole('checkbox', { name: /delete all associated parts/i })
    ).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /yes/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /cancel/i })).toBeInTheDocument()
  })

  it('should delete only the section by default', async () => {
    const user = userEvent.setup()

    const onClose = vi.fn()
    const songId = 'song-1'

    let searchParams: URLSearchParams
    server.use(
      http.delete(`/songs/sections/${section.id}/from/${songId}`, ({ request }) => {
        searchParams = new URL(request.url).searchParams
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    reduxRender(
      withToastify(
        <DeleteSongSectionModal opened={true} onClose={onClose} section={section} songId={songId} />
      )
    )

    await user.click(screen.getByRole('button', { name: /yes/i }))

    expect(searchParams.get('withParts')).toBe('false')
    expect(screen.getByText(`${section.name} deleted!`)).toBeInTheDocument()
    expect(onClose).toHaveBeenCalledOnce()
  })

  it('should delete the section and associated parts when selected', async () => {
    const user = userEvent.setup()

    const songId = 'song-1'

    let searchParams: URLSearchParams
    server.use(
      http.delete(`/songs/sections/${section.id}/from/${songId}`, ({ request }) => {
        searchParams = new URL(request.url).searchParams
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    reduxRender(
      withToastify(
        <DeleteSongSectionModal opened={true} onClose={vi.fn()} section={section} songId={songId} />
      )
    )

    await user.click(screen.getByRole('checkbox', { name: /delete all associated parts/i }))
    await user.click(screen.getByRole('button', { name: /yes/i }))

    expect(searchParams.get('withParts')).toBe('true')
  })
})
