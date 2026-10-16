import {
  defaultDropAnimationSideEffects,
  DndContext,
  DragOverlay,
  useDroppable,
  closestCenter,
  useSensor,
  useSensors,
  KeyboardSensor,
  PointerSensor,
} from '@dnd-kit/core';
import { CSS } from '@dnd-kit/utilities';

import {
  arrayMove,
  SortableContext,
  sortableKeyboardCoordinates,
  useSortable,
  verticalListSortingStrategy,
} from '@dnd-kit/sortable';
import {
  AvatarGroup,
  Box,
  Button,
  Chip,
  IconButton,
  LinearProgress,
  Menu,
  MenuItem,
  Paper,
  Stack,
  Typography,
} from '@mui/material';
import { useEffect, useMemo, useState } from 'react';

import ModalAddNewTask from './Modals/ModalAddNewTask';
import ModalDetailTask from './Modals/ModalDetailTask';

import SidebarLayout from '@/components/layouts/SidebarLayout';
import EditProject from './EditProject';
import { useParams } from 'react-router';
import boards from '@/services/api/boards';
import datetime from '@/utils/datetime';
import lists from '@/services/api/lists';
import cards from '@/services/api/cards';
import Avatar from '@/components/ui/Avatar';

const ListDroppable = ({ id, data, children }) => {
  const { setNodeRef } = useDroppable({
    id,
    data,
  });

  return (
    <Box
      sx={{
        width: '100%',
      }}
      ref={setNodeRef}
    >
      {children}
    </Box>
  );
};

const CardItem = ({ id, data, onClick }) => {
  const {
    setNodeRef,
    attributes,
    listeners,
    transform,
    transition,
    isDragging,
  } = useSortable({
    id,
    data,
  });

  const style = {
    opacity: isDragging ? 0.5 : 1,
    transform: CSS.Translate.toString(transform),
    transition,
    padding: 1,
    border: '1px solid #ccc',
    borderRadius: '4px',
    backgroundColor: '#f9f9f9',
    ':hover': {
      cursor: 'pointer',
      border: '1px solid #1976d2',
    },
    minHeight: '5rem',
  };

  const handleClick = () => {
    onClick(data);
  };

  return (
    <Box
      {...attributes}
      {...listeners}
      onClick={handleClick}
      sx={style}
      ref={setNodeRef}
    >
      <Typography variant="body1" fontWeight={'bold'}>
        {data.title}
      </Typography>
      <Box
        mt={1}
        sx={{
          maxWidth: '100%',
          // maxHeight: '3.5rem',
          overflow: 'hidden',
        }}
      >
        <Typography variant="caption">{data.description}</Typography>
      </Box>
    </Box>
  );
};

const DetailProject = () => {
  const { id } = useParams();

  const [toggleModalAddNewTask, setToggleModalAddNewTask] = useState(false);
  const [toggleModalDetailTask, setToggleModalDetailTask] = useState(false);
  const [modalCardData, setModalCardData] = useState(null);

  const [active, setActive] = useState(null);
  const [over, setOver] = useState(null);

  const [boardData, setBoardData] = useState({});
  const [isLoading, setLoading] = useState(false);
  const [isLoadingBoardLists, setLoadingBoardLists] = useState(false);

  const [boardList, setBoardList] = useState([]);
  // const [boardListCards, setBoardListCards] = useState([]);

  const [boardItems, setBoardItems] = useState([]);

  // fetch board detail
  useEffect(() => {
    const fetchBoard = async () => {
      setLoading(true);
      const response = await boards.detail(id);

      setBoardData(response.data.data);
      setLoading(false);
    };
    fetchBoard();
  }, [id]);

  // fetch lists on board
  useEffect(() => {
    const fetchBoardLists = async () => {
      setLoadingBoardLists(true);
      const response = await boards.lists(id);
      const lists = response.data.data.Lists;
      const positions = response.data.data.Positions;
      const boardListData = positions.map((boardId, idx) => {
        if (boardId === lists[idx].public_id) {
          return {
            id: boardId,
            title: lists[idx].title,
          };
        }
        return {};
      });

      setBoardList(boardListData);
      setLoadingBoardLists(false);
    };
    fetchBoardLists();
  }, [id]);

  // fetch cards on list
  useEffect(() => {
    if (boardList && boardList.length > 0) {
      fetchCardsOnList();
    }
  }, [boardList]);

  const updatePositions = async (boardItems) => {
    const groupedMap = new Map();

    boardItems.forEach((item) => {
      const { boardListId, id } = item;

      // Cek apakah boardListId sudah ada di Map
      if (groupedMap.has(boardListId)) {
        // Jika sudah ada, tambahkan id ke array yang ada
        groupedMap.get(boardListId).positions.push(id);
      } else {
        // Jika belum ada, buat entri baru di Map
        groupedMap.set(boardListId, {
          boardListId: boardListId,
          positions: [id],
        });
      }
    });

    // Ubah Map values menjadi array
    const cardsGroup = Array.from(groupedMap.values());

    // console.log(grouped);

    for (let cardItem of cardsGroup) {
      await lists.updatePosition(cardItem.boardListId, {
        positions: cardItem.positions,
      });
    }
  };

  const updateCardList = async (active, over) => {
    // console.log({ active, over });
    // console.log({
    //   list_id: over.id,
    //   title: active.title,
    //   description: active.description,
    //   due_date: active.due_date,
    //   position: getItemsByBoardId(over.id).length + 1,
    // });
    await cards.update(active.id, {
      list_id: over.id ?? over.boardListId,
      title: active.title,
      description: active.description,
      due_date: active.due_date,
      position: getItemsByBoardId(over.id).length + 1,
    });
  };

  const fetchCardsOnList = async () => {
    // loop inside lists
    const cardsData = [];
    for (let idxList = 0; idxList < boardList.length; idxList++) {
      const cardsOnListData = await lists.getCardsOnList(boardList[idxList].id);

      for (
        let idxCard = 0;
        idxCard < cardsOnListData.data.data.length;
        idxCard++
      ) {
        const _cardsData = cardsOnListData.data.data[idxCard];
        const isCardExist = cardsData.find(
          (item) => item.id === _cardsData.public_id,
        );
        if (!isCardExist) {
          cardsData.push({
            id: _cardsData.public_id,
            title: _cardsData.title,
            description: _cardsData.description,
            boardListId: boardList[idxList].id,
            position: _cardsData.position,
            due_date: _cardsData.due_date,
          });
        }
      }
    }
    setBoardItems(cardsData);
  };

  const getItemsByBoardId = (boardId) => {
    const filtered = boardItems.filter((item) => item.boardListId === boardId);
    return filtered;
  };

  const handleAddNewCard = (board) => {
    setModalCardData(board);
    setToggleModalAddNewTask(true);
  };

  const sensors = useSensors(
    useSensor(PointerSensor, {
      activationConstraint: {
        distance: 5,
      },
    }),
    useSensor(KeyboardSensor, {
      coordinateGetter: sortableKeyboardCoordinates,
    }),
  );

  const handleSelectedCardItem = (data) => {
    setToggleModalDetailTask(true);
    setModalCardData(data);
  };

  const handleSaveAddNewTask = async (data) => {
    await cards.create({
      ...data,
      position: getItemsByBoardId(data.list_id).length + 1,
    });

    await fetchCardsOnList();

    setToggleModalAddNewTask(false);
    setModalCardData(null);
  };

  return (
    <SidebarLayout
      pageTitle={boardData?.title}
      breadcrumbs={[
        {
          label: 'Daftar Proyek',
          href: '/projects',
        },
        {
          label: boardData?.title,
        },
      ]}
    >
      <Box
        sx={{
          mb: 2,
        }}
      >
        <Stack spacing={1} direction={'row'} justifyContent={'space-between'}>
          <Stack spacing={1} direction={'row'} alignItems={'center'}>
            <Chip
              label={`Deadline: ${datetime.format(boardData?.due_date)}`}
              color="error"
            />
            <AvatarGroup>
              <Avatar text={'Muhammad Agung'} />
              <Avatar text={'Agung Rizkyana'} />
              <Avatar text={'Triadynata Zamhur'} />
            </AvatarGroup>
          </Stack>
          <Box>
            <EditProject data={boardData} />
          </Box>
        </Stack>
      </Box>
      <Box
        sx={{
          p: 1,
        }}
      >
        {isLoading && <LinearProgress />}
        <DndContext
          sensors={sensors}
          collisionDetection={closestCenter}
          onDragStart={(event) => {
            setActive(event.active.data.current);
          }}
          onDragOver={(event) => {
            setOver(event.over?.data.current);
          }}
          onDragEnd={async (event) => {
            const { active, over } = event;
            if (!over) return;

            if (
              active.data.current.boardListId !== over.data.current.boardListId
            ) {
              const updatedBoardItems = [...boardItems].map((item) => {
                if (item.id === active.data.current.id) {
                  return {
                    ...item,
                    boardListId:
                      over.data.current.boardListId ?? over.data.current.id,
                  };
                }
                return {
                  ...item,
                };
              });
              setBoardItems(updatedBoardItems);
              await updateCardList(active.data.current, over.data.current);
              return;
            }

            if (active.id !== over.id) {
              const oldIndex = [...boardItems].findIndex(
                (item) => item.id === active.id,
              );

              const newIndex = [...boardItems].findIndex(
                (item) => item.id === over.id,
              );

              const updatedItems = arrayMove(boardItems, oldIndex, newIndex);
              setBoardItems(updatedItems);
              await updatePositions(updatedItems);
              return;
            }

            setActive(null);
            setOver(null);
          }}
          onDragCancel={() => setActive(null)}
        >
          <Box
            sx={{
              overflow: 'hidden',
            }}
          >
            <Box
              sx={{
                display: 'flex',
                flexDirection: 'row',
                gap: 3,
              }}
            >
              {boardList.map((board, idx) => (
                <ListDroppable key={board.id} id={board.id} data={board}>
                  <SortableContext
                    key={board.id}
                    items={getItemsByBoardId(board.id)}
                    strategy={verticalListSortingStrategy}
                  >
                    <Box
                      sx={{
                        background: '#fff',
                        width: '100%',
                        border: `1px solid ${
                          over?.boardListId === board.id ||
                          over?.id === board.id
                            ? '#1976d2'
                            : '#ccc'
                        }`,
                        borderRadius: 1,
                        ':hover': {
                          cursor: 'pointer',
                        },
                      }}
                    >
                      <Box
                        sx={{
                          display: 'flex',
                          justifyContent: 'space-between',
                          alignItems: 'center',
                          padding: 2,
                          borderBottom: '1px solid #ccc',
                          background:
                            boardList.length - 1 === idx ? '#aff2a7' : '#fff',
                          borderTopRightRadius: '4px',
                          borderTopLeftRadius: '4px',
                        }}
                      >
                        <Typography variant="h6">{board.title}</Typography>
                        <Typography variant="h6">
                          {getItemsByBoardId(board.id).length}
                        </Typography>
                      </Box>
                      <Box
                        sx={{
                          display: 'flex',
                          flexDirection: 'column',
                          gap: 1,
                          maxHeight: 'calc(100vh - 6rem)',
                          overflowY: 'auto',
                          width: '100%',
                          p: 2,
                        }}
                      >
                        {getItemsByBoardId(board.id).map((item) => (
                          <CardItem
                            key={item.id}
                            id={item.id}
                            data={item}
                            onClick={handleSelectedCardItem}
                          />
                        ))}
                      </Box>
                      {idx === 0 && (
                        <Box sx={{ textAlign: 'center', p: 2 }}>
                          <Button
                            type="button"
                            variant="contained"
                            size="large"
                            onClick={() => handleAddNewCard(board)}
                            fullWidth
                          >
                            Tambah Tugas Baru
                          </Button>
                        </Box>
                      )}
                    </Box>
                  </SortableContext>
                </ListDroppable>
              ))}
            </Box>

            <DragOverlay
              dropAnimation={{
                sideEffects: defaultDropAnimationSideEffects({
                  styles: {
                    active: {
                      opacity: '0.4',
                    },
                  },
                }),
              }}
            >
              {active ? <CardItem id={active.id} data={active} /> : null}
            </DragOverlay>
          </Box>
        </DndContext>
      </Box>
      <ModalAddNewTask
        open={toggleModalAddNewTask}
        handleClose={() => {
          setToggleModalAddNewTask(false);
          setModalCardData(null);
        }}
        onSaved={handleSaveAddNewTask}
        data={modalCardData}
      />
      <ModalDetailTask
        open={toggleModalDetailTask}
        data={modalCardData}
        handleClose={() => {
          setToggleModalDetailTask(false);
          setModalCardData(null);
        }}
      />
    </SidebarLayout>
  );
};

export default DetailProject;
