from __future__ import annotations

import math
from typing import Iterable

import numpy as np


def poisson_pmf(k: int, lam: float) -> float:
    if lam <= 0:
        return 1.0 if k == 0 else 0.0
    return math.exp(-lam) * (lam**k) / math.factorial(k)


def poisson_score_matrix(lam_home: float, lam_away: float, max_goals: int) -> np.ndarray:
    max_goals = max(0, int(max_goals))
    ph = np.array([poisson_pmf(i, lam_home) for i in range(max_goals + 1)], dtype=float)
    pa = np.array([poisson_pmf(i, lam_away) for i in range(max_goals + 1)], dtype=float)
    return np.outer(ph, pa)


def cosine_similarity(a: np.ndarray, b: np.ndarray) -> float:
    na = np.linalg.norm(a)
    nb = np.linalg.norm(b)
    if na == 0.0 or nb == 0.0:
        return 0.0
    return float(np.dot(a, b) / (na * nb))


def zscore_normalize(matrix: np.ndarray) -> tuple[np.ndarray, np.ndarray, np.ndarray]:
    mean = matrix.mean(axis=0)
    std = matrix.std(axis=0)
    std = np.where(std == 0.0, 1.0, std)
    return (matrix - mean) / std, mean, std


def ewma_weights(n: int, half_life: float) -> np.ndarray:
    n = int(n)
    if n <= 0:
        return np.array([], dtype=float)
    hl = float(half_life) if half_life and half_life > 0 else 5.0
    # Newest match has age=0. Oldest has age=n-1.
    ages = np.arange(n, dtype=float)
    decay = math.log(2.0) / hl
    w = np.exp(-decay * ages)
    return w / w.sum()


def weighted_avg(values: Iterable[float], weights: np.ndarray) -> float:
    v = np.array(list(values), dtype=float)
    if v.size == 0:
        return 0.0
    w = weights[: v.size]
    if w.size != v.size:
        # Fallback: uniform.
        return float(v.mean())
    return float(np.dot(v, w))

