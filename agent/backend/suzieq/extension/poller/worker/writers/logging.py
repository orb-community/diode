"""
This module contains the logic of the writer for the 'logging' mode
"""
import logging

from suzieq.poller.worker.writers.output_worker import OutputWorker

logger = logging.getLogger(__name__)

class LoggingOutputWorker(OutputWorker):
    """LoggingOutputWorker is used to write poller output as logs
    """
    def __init__(self, **kwargs):
        self.data_directory = kwargs.get('data_dir')

    def write_data(self, data):
        """Write the output of the commands into stdout

        Args:
            data (Dict): dictionary containing the data to store.
        """
        if not data["records"]:
            return

        for record in data["records"]:
            logger.warning(f'topic: {data["topic"]}, record: {record}')


