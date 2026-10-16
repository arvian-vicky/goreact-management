import network from '@/utils/network';

const lists = {
  async create(payload) {
    return network.post(`/api/v1/lists`, payload);
  },
  async updatePosition(id, payload) {
    return network.put(`/api/v1/lists/${id}/positions`, payload);
  },
  async getCardsOnList(id) {
    return network.get(`/api/v1/lists/${id}/cards`);
  },
};

export default lists;
