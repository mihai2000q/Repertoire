import { http, HttpResponse } from 'msw'
import {
  emptySongPart,
  emptySongSection,
  reduxRender,
  withToastify
} from '../../../../../../test-utils.tsx'
import { screen } from '@testing-library/react'
import { setupServer } from 'msw/node'
import { userEvent } from '@testing-library/user-event'
import DeleteSongSectionsModal from './DeleteSongSectionsModal.tsx'
import { BulkDeleteSongSectionsRequest } from '../../types/requests/SongSectionRequests.ts'
import { SongSection } from '../../../../../../types/models/Song.ts'

describe('Delete Song Sections Modal', () => {
  const sections: SongSection[] = [
    {
      ...emptySongSection,
      id: 'section-1',
      name: 'Verse',
      parts: [{ ...emptySongPart, id: 'part-1' }]
    },
    {
      ...emptySongSection,
      id: 'section-2',
      name: 'Chorus',
      parts: [
        { ...emptySongPart, id: 'part-1' },
        { ...emptySongPart, id: 'part-2' }
      ]
    }
  ]

  const server = setupServer()

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  it('should render', async () => {
    reduxRender(
      <DeleteSongSectionsModal
        sections={sections}
        songId={'song-1'}
        opened={true}
        onClose={vi.fn()}
      />
    )

    expect(screen.getByRole('dialog', { name: /delete sections/i })).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: /delete sections/i })).toBeInTheDocument()
    expect(screen.getByText(/are you sure/i)).toBeInTheDocument()
    expect(screen.getByRole('checkbox', { name: /delete all associated parts/i })).toBeEnabled()
    expect(screen.getByText(/2 parts/i)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /yes/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /cancel/i })).toBeInTheDocument()
  })

  it('should close when cancelling', async () => {
    const user = userEvent.setup()
    const onClose = vi.fn()

    reduxRender(
      <DeleteSongSectionsModal
        sections={sections}
        songId={'song-1'}
        opened={true}
        onClose={onClose}
      />
    )

    await user.click(screen.getByRole('button', { name: /cancel/i }))

    expect(onClose).toHaveBeenCalledOnce()
  })

  it('should disable deleting associated parts when no parts exist', () => {
    reduxRender(
      <DeleteSongSectionsModal
        sections={[{ ...emptySongSection, id: 'section-1' }]}
        songId={'song-1'}
        opened={true}
        onClose={vi.fn()}
      />
    )

    expect(
      screen.getByRole('checkbox', { name: /delete all associated parts \(0 parts\)/i })
    ).toBeDisabled()
  })

  it('should include the deduplicated associated part count when selected', async () => {
    const user = userEvent.setup()

    server.use(
      http.put(`/songs/sections/bulk-delete`, () => HttpResponse.json({ message: 'it worked' }))
    )

    reduxRender(
      withToastify(
        <DeleteSongSectionsModal
          sections={sections}
          songId={'song-1'}
          opened={true}
          onClose={vi.fn()}
        />
      )
    )

    const checkbox = screen.getByRole('checkbox', { name: /delete all associated parts/i })
    await user.click(checkbox)
    expect(checkbox).toBeChecked()

    await user.click(screen.getByRole('button', { name: /yes/i }))

    expect(await screen.findByText('2 sections deleted with 2 parts!')).toBeInTheDocument()
  })

  it('should delete song sections in bulk and call onDelete', async () => {
    const user = userEvent.setup()

    let capturedRequest: BulkDeleteSongSectionsRequest | undefined
    server.use(
      http.put(`/songs/sections/bulk-delete`, async (req) => {
        capturedRequest = (await req.request.json()) as BulkDeleteSongSectionsRequest
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    const songId = 'song-1'
    const onClose = vi.fn()
    const onDelete = vi.fn()

    reduxRender(
      withToastify(
        <DeleteSongSectionsModal
          sections={sections.map((s) => ({ ...s, parts: [] }))}
          songId={songId}
          opened={true}
          onClose={onClose}
          onDelete={onDelete}
        />
      )
    )

    await user.click(screen.getByRole('button', { name: /yes/i }))

    expect(await screen.findByText(`${sections.length} sections deleted!`)).toBeInTheDocument()
    expect(onClose).toHaveBeenCalledOnce()
    expect(onDelete).toHaveBeenCalledOnce()
    expect(capturedRequest).toStrictEqual({
      ids: sections.map((s) => s.id),
      partIds: [],
      songId: songId
    })
  })

  it('should delete song sections in bulk and call onDelete with parts', async () => {
    const user = userEvent.setup()

    let capturedRequest: BulkDeleteSongSectionsRequest | undefined
    server.use(
      http.put(`/songs/sections/bulk-delete`, async (req) => {
        capturedRequest = (await req.request.json()) as BulkDeleteSongSectionsRequest
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    const songId = 'song-1'
    const onClose = vi.fn()
    const onDelete = vi.fn()

    const expectedIds: string[] = []
    const expectedPartIds: Set<string> = new Set<string>()
    sections.forEach((section) => {
      expectedIds.push(section.id)
      section.parts.forEach((part) => {
        expectedPartIds.add(part.id)
      })
    })

    reduxRender(
      withToastify(
        <DeleteSongSectionsModal
          sections={sections}
          songId={songId}
          opened={true}
          onClose={onClose}
          onDelete={onDelete}
        />
      )
    )

    await user.click(screen.getByRole('checkbox', { name: /delete all associated parts/i }))
    await user.click(screen.getByRole('button', { name: /yes/i }))

    expect(
      await screen.findByText(
        `${sections.length} sections deleted with ${expectedPartIds.size} parts!`
      )
    ).toBeInTheDocument()
    expect(onClose).toHaveBeenCalledOnce()
    expect(onDelete).toHaveBeenCalledOnce()
    expect(capturedRequest).toStrictEqual({
      ids: expectedIds,
      partIds: Array.from(expectedPartIds),
      songId: songId
    })
  })
})
