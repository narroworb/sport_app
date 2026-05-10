from clickhouse_driver import Client

class ClickHouseClientWrapper:
    """Адаптер для clickhouse-driver, чтобы имитировать интерфейс clickhouse-connect"""
    
    def __init__(self, host, port, username, password, database):
        self._client = Client(
            host=host,
            port=port,
            user=username,
            password=password,
            database=database
        )
    
    def execute(self, query, parameters=None):
        """Основной метод для выполнения запросов"""
        print(f"Executing query: {query[:100]}...")  # Для отладки
        if parameters:
            return self._client.execute(query, parameters)
        return self._client.execute(query)

    def query(self, query, parameters=None):
        """Имитация метода query из clickhouse-connect"""
        if parameters:
            # Преобразуем параметры в формат clickhouse-driver
            result = self._client.execute(query, parameters)
        else:
            result = self._client.execute(query)
        
        # Возвращаем объект, похожий на результат clickhouse-connect
        return ClickHouseResultWrapper(result)
    
    def query_df(self, query):
        """Если нужен DataFrame"""
        import pandas as pd
        result = self._client.execute(query)
        return pd.DataFrame(result)
    
    def command(self, query):
        """Для команд (INSERT, CREATE и т.д.)"""
        return self._client.execute(query)

class ClickHouseResultWrapper:
    """Адаптер для результатов запроса"""
    
    def __init__(self, result):
        self._result = result
        self.result_rows = result
    
    def __iter__(self):
        return iter(self._result)
    
    def __getitem__(self, idx):
        return self._result[idx]
    
    @property
    def rows(self):
        return self._result
    
    def first_row(self):
        if self._result:
            return self._result[0]
        return None
    

def create_client(database_addr, database, username, password):
    """Создает клиент ClickHouse с единым интерфейсом"""
    # Парсим адрес: "clickhouse:9000" или "http://clickhouse:9000"
    addr = database_addr.replace('http://', '').replace('https://', '')
    
    if ':' in addr:
        host, port = addr.split(':')
        port = int(port)
    else:
        host = addr
        port = 9000  # порт по умолчанию для native protocol
    
    return ClickHouseClientWrapper(
        host=host,
        port=port,
        username=username,
        password=password,
        database=database
    )