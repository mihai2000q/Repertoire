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
import DeleteSongSectionsAndPartsModal from './DeleteSongSectionsAndPartsModal.tsx'
import { BulkDeleteSongSectionsRequest } from '../../types/requests/SongSectionRequests.ts'
import { SongSection } from '../../../../../../types/models/Song.ts'

describe('Delete Song Sections And Parts Modal', () => {
  const sections: SongSection[] = [
    {
      ...emptySongSection,
      id: 'section-1',
      name: 'Verse 123',
      rehearsals: 12,
      songSectionType: { id: '1', name: 'Verse' },
      parts: [{ ...emptySongPart, id: 'part-1', name: 'Shared part', rehearsals: 1 }]
    },
    {
      ...emptySongSection,
      id: 'section-2',
      name: 'Chorus 123',
      rehearsals: 6,
      songSectionType: { id: '1', name: 'Chorus' },
      parts: [
        { ...emptySongPart, id: 'part-1', name: 'Shared part', rehearsals: 1 },
        { ...emptySongPart, id: 'part-2', name: 'Unique part', rehearsals: 2 },
        { ...emptySongPart, id: 'part-3', name: 'Not Included', rehearsals: 3 }
      ]
    }
  ]
  const sectionParts = [
    { ...emptySongPart, id: 'part-1' },
    { ...emptySongPart, id: 'part-1' },
    { ...emptySongPart, id: 'part-2' }
  ]

  const server = setupServer()

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  it('should render sections and its selected parts', async () => {
    const user = userEvent.setup()

    const seen: Set<string> = new Set<string>()
    const duplicateParts: Set<string> = new Set<string>()
    sectionParts.forEach((p) => {
      if (seen.has(p.id)) {
        duplicateParts.add(p.id)
      }
      seen.add(p.id)
    })

    reduxRender(
      <DeleteSongSectionsAndPartsModal
        sections={sections}
        sectionParts={sectionParts}
        songId={'song-1'}
        opened={true}
        onClose={vi.fn()}
      />
    )

    expect(
      screen.getByRole('dialog', { name: /delete song sections with parts/i })
    ).toBeInTheDocument()
    expect(screen.getByText(/are you sure/i)).toBeInTheDocument()

    for (const section of sections) {
      expect(screen.getByText(section.name)).toBeInTheDocument()
      expect(screen.getByText(section.songSectionType.name)).toBeInTheDocument()
      expect(screen.getByText(section.rehearsals)).toBeInTheDocument()

      for (const part of section.parts) {
        if (!sectionParts.some((p) => p.id === part.id)) continue
        if (duplicateParts.has(part.id)) {
          const duplicates = sectionParts.filter((p) => p.id === part.id)
          expect(screen.getAllByText(part.name)).toHaveLength(duplicates.length)
          expect(screen.getAllByText(part.rehearsals)).toHaveLength(duplicates.length)
          expect(screen.getAllByLabelText(`${part.name}-duplicate`)).toHaveLength(duplicates.length)
          await user.hover(screen.getAllByLabelText(`${part.name}-duplicate`)[0])
          expect(await screen.findByText('This is a duplicate')).toBeInTheDocument()
        } else {
          expect(screen.getByText(part.name)).toBeInTheDocument()
          expect(screen.getByText(part.rehearsals)).toBeInTheDocument()
        }
      }
    }

    expect(screen.getByRole('button', { name: /cancel/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /yes/i })).toBeInTheDocument()
  })

  it('should delete sections and parts in bulk', async () => {
    const user = userEvent.setup()

    let capturedRequest: BulkDeleteSongSectionsRequest
    server.use(
      http.put(`/songs/sections/bulk-delete`, async (request) => {
        capturedRequest = (await request.request.json()) as BulkDeleteSongSectionsRequest
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    const songId = 'song-1'
    const onDelete = vi.fn()
    const onClose = vi.fn()

    const partIds: Set<string> = new Set<string>()
    sectionParts.forEach((p) => {
      partIds.add(p.id)
    })

    reduxRender(
      withToastify(
        <DeleteSongSectionsAndPartsModal
          sections={sections}
          sectionParts={sectionParts}
          songId={songId}
          opened={true}
          onClose={onClose}
          onDelete={onDelete}
        />
      )
    )

    await user.click(screen.getByRole('button', { name: /yes/i }))

    expect(await screen.findByText('2 sections deleted with 2 parts!')).toBeInTheDocument()
    expect(capturedRequest).toStrictEqual({
      ids: sections.map((s) => s.id),
      partIds: Array.from(partIds),
      songId: songId
    })
    expect(onDelete).toHaveBeenCalledOnce()
    expect(onClose).toHaveBeenCalledOnce()
  })

  it('should close when cancelling', async () => {
    const user = userEvent.setup()

    const onClose = vi.fn()

    reduxRender(
      <DeleteSongSectionsAndPartsModal
        sections={sections}
        sectionParts={sectionParts}
        songId={'song-1'}
        opened={true}
        onClose={onClose}
      />
    )

    await user.click(screen.getByRole('button', { name: /cancel/i }))

    expect(onClose).toHaveBeenCalledOnce()
  })
})
