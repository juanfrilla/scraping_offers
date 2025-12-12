from datetime import datetime
from typing import Optional

from pydantic import BaseModel, HttpUrl


class SalaryModel(BaseModel):
    min: Optional[int]
    max: Optional[int]
    currency: Optional[str]


class JobOfferModel(BaseModel):
    external_id: str
    title: str
    description: str
    city: Optional[str]
    link: Optional[HttpUrl]
    salary: Optional[SalaryModel]
    modality: Optional[str]
    published_at: Optional[datetime]
    company_name: Optional[str]
    platform: Optional[str]
