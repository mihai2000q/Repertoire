import { emptySongSection, reduxRender, withToastify } from '../../../../../../../../test-utils.tsx'
import { SongSection, SongSectionType } from '../../../../../../../../types/models/Song.ts'
import { setupServer } from 'msw/node'
import { screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import { http, HttpResponse } from 'msw'
import { UpdateSongSectionRequest } from '../../types/requests/SongSectionRequests.ts'
import EditSongSectionModal from './EditSongSectionModal.tsx'

describe('Edit Song Description Modal', () => {
  const sectionTypes: SongSectionType[] = [
    {
      id: '1',
      name: 'Solo'
    },
    {
      id: '2',
      name: 'Riff'
    }
  ]

  const section: SongSection = {
    ...emptySongSection,
    id: 'some-id',
    name: 'section 1',
    songSectionType: sectionTypes[1]
  }

  const handlers = [
    http.get(`/songs/sections/types`, () => {
      return HttpResponse.json(sectionTypes)
    })
  ]

  const server = setupServer(...handlers)

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  it('should render', async () => {
    const user = userEvent.setup()

    reduxRender(<EditSongSectionModal opened={true} onClose={() => {}} section={section} />)

    expect(screen.getByRole('dialog', { name: /edit song section/i })).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: /edit song section/i })).toBeInTheDocument()

    expect(screen.getByRole('textbox', { name: /name/i })).toBeInTheDocument()
    expect(screen.getByRole('textbox', { name: /name/i })).not.toBeInvalid()
    expect(screen.getByRole('textbox', { name: /name/i })).toHaveValue(section.name)

    expect(await screen.findByRole('combobox', { name: /type/i })).toBeInTheDocument()
    expect(screen.getByRole('combobox', { name: /type/i })).not.toBeInvalid()
    expect(await screen.findByRole('combobox', { name: /type/i })).toHaveValue(
      section.songSectionType.name
    )

    expect(screen.getByRole('button', { name: /save/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /save/i })).toBeDisabled()
    await user.hover(screen.getByRole('button', { name: /save/i }))
    expect(await screen.findByText(/need to make a change/i)).toBeInTheDocument()
  })

  it('should send update request when the field values have changed', async () => {
    const user = userEvent.setup()

    const newName = 'New Section Name'
    const newType = sectionTypes[0]
    const onClose = vitest.fn()

    let capturedRequest: UpdateSongSectionRequest
    server.use(
      http.put('/songs/sections', async (req) => {
        capturedRequest = (await req.request.json()) as UpdateSongSectionRequest
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    const [{ rerender }] = reduxRender(
      withToastify(<EditSongSectionModal opened={true} onClose={onClose} section={section} />)
    )

    const nameField = screen.getByRole('textbox', { name: /name/i })
    const typeField = screen.getByRole('combobox', { name: /type/i })
    const saveButton = screen.getByRole('button', { name: /save/i })

    await user.clear(nameField)
    await user.type(nameField, newName)

    await user.click(typeField)
    await user.click(await screen.findByText(newType.name))

    expect(saveButton).not.toBeDisabled()
    await user.click(screen.getByRole('button', { name: /save/i }))

    expect(capturedRequest).toStrictEqual({
      id: section.id,
      name: newName,
      typeId: newType.id,
      partIds: []
    })
    expect(onClose).toHaveBeenCalledOnce()

    expect(await screen.findByText(`${newName} updated!`)).toBeInTheDocument()

    rerender(
      <EditSongSectionModal
        opened={true}
        onClose={onClose}
        section={{
          ...section,
          name: newName,
          songSectionType: newType
        }}
      />
    )
    expect(screen.getByRole('button', { name: /save/i })).toBeDisabled()
  })

  it('should keep the save button disabled when the field values have not changed', async () => {
    const user = userEvent.setup()

    reduxRender(<EditSongSectionModal opened={true} onClose={() => {}} section={section} />)

    const nameField = screen.getByRole('textbox', { name: /name/i })
    const typeField = screen.getByRole('combobox', { name: /type/i })
    const saveButton = screen.getByRole('button', { name: /save/i })

    // change name
    await user.clear(nameField)
    await user.type(nameField, section.name + '1')
    expect(saveButton).not.toBeDisabled()

    // reset name
    await user.clear(nameField)
    await user.type(nameField, section.name)
    expect(saveButton).toBeDisabled()

    // change type
    await user.click(typeField)
    await user.click(await screen.findByText(sectionTypes[0].name))
    expect(saveButton).not.toBeDisabled()

    // reset type
    await user.click(typeField)
    await user.click(await screen.findByText(section.songSectionType.name))
    expect(saveButton).toBeDisabled()
  })

  it('should validate fields', async () => {
    const user = userEvent.setup()

    reduxRender(<EditSongSectionModal opened={true} onClose={() => {}} section={section} />)

    const nameField = screen.getByRole('textbox', { name: /name/i })

    // invalidate name
    await user.clear(nameField)
    expect(nameField).toBeInvalid()
  })
})
